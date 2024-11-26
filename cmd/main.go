package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"AI-RSS-Telegram-Bot/internal/summary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"

	"AI-RSS-Telegram-Bot/internal/bot"
	"AI-RSS-Telegram-Bot/internal/bot/middleware"
	"AI-RSS-Telegram-Bot/internal/botkit"
	"AI-RSS-Telegram-Bot/internal/config"
	"AI-RSS-Telegram-Bot/internal/fetcher"
	"AI-RSS-Telegram-Bot/internal/notifier"
	"AI-RSS-Telegram-Bot/internal/storage"
)

func main() {
	botAPI, err := tgbotapi.NewBotAPI(config.Get().TelegramBotToken)
	if err != nil {
		log.Printf("[ERROR] failed to create botAPI: %v", err)
		return
	}

	db, err := sql.Open("postgres", config.Get().DatabaseDSN)
	if err != nil {
		log.Printf("[ERROR] failed to open db connection: %v", err)
		return
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("[ERROR] failed to close db: %v", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Printf("[ERROR] failed to ping db: %v", err)
		return
	}

	log.Println("[INFO] Setting bot commands...")
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Запустить бота"},
		{Command: "addsource", Description: "Добавить новый источник"},
		{Command: "setpriority", Description: "Установить приоритет источника"},
		{Command: "getsource", Description: "Получить информацию об источнике"},
		{Command: "listsources", Description: "Показать список всех источников"},
		{Command: "deletesource", Description: "Удалить источник"},
	}

	_, err = botAPI.Request(tgbotapi.NewSetMyCommands(commands...))
	if err != nil {
		log.Fatalf("[ERROR] failed to set bot commands: %v", err)
	}
	log.Println("[INFO] Bot commands set successfully!")

	proxyURL, err := url.Parse(config.Get().ProxyURL)
	if err != nil {
		log.Printf("[ERROR] failed to parse proxy URL: %v", err)
		return
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	var (
		articleStorage = storage.NewArticleStorage(db)
		sourceStorage  = storage.NewSourceStorage(db)
		fetcherService = fetcher.New(
			articleStorage,
			sourceStorage,
			config.Get().FetchInterval,
			config.Get().FilterKeywords,
		)

		summarizer = summary.NewOpenAISummarizer(
			config.Get().OpenAIKey,
			config.Get().OpenAIModel,
			config.Get().OpenAIPrompt,
			httpClient,
		)
		notifierService = notifier.New(
			articleStorage,
			summarizer,
			botAPI,
			config.Get().NotificationInterval,
			2*config.Get().FetchInterval,
			config.Get().TelegramChannelID,
		)
	)

	newsBot := botkit.New(botAPI)

	newsBot.RegisterCmdView(
		"start",
		bot.ViewCmdStart(),
	)
	newsBot.RegisterCmdView(
		"addsource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdAddSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"setpriority",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdSetPriority(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"getsource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdGetSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"listsources",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdListSource(sourceStorage),
		),
	)
	newsBot.RegisterCmdView(
		"deletesource",
		middleware.AdminsOnly(
			config.Get().TelegramChannelID,
			bot.ViewCmdDeleteSource(sourceStorage),
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func(ctx context.Context) {
		if err := fetcherService.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to run fetcher: %v", err)
				return
			}
			log.Printf("[INFO] fetcher stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := notifierService.Start(ctx); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to run notifier: %v", err)
				return
			}
			log.Printf("[INFO] notifier stopped")
		}
	}(ctx)

	go func(ctx context.Context) {
		if err := http.ListenAndServe("0.0.0.0:8080", mux); err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Printf("[ERROR] failed to run http server: %v", err)
				return
			}
			log.Printf("[INFO] http server stopped")
		}
	}(ctx)

	if err := newsBot.Run(ctx); err != nil {
		log.Printf("[ERROR] failed to run botkit: %v", err)
	}
}
