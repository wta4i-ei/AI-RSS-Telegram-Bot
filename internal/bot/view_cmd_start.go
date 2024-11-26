package bot

import (
	"AI-RSS-Telegram-Bot/internal/botkit"
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ViewCmdStart() botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		// Приветственное сообщение
		welcomeMessage := "Привет! Я твой бот для работы с источниками. Чем могу помочь?"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMessage)
		if _, err := bot.Send(msg); err != nil {
			return fmt.Errorf("ошибка отправки приветственного сообщения %w", err)
		}
		return nil
	}
}
