package bot

import (
	"context"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"AI-RSS-Telegram-Bot/internal/botkit"
	"AI-RSS-Telegram-Bot/internal/model"
)

type SourceStorage interface {
	Add(ctx context.Context, source model.Source) (int64, error)
}

func ViewCmdAddSource(storage SourceStorage) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args := update.Message.CommandArguments()

		var name, url string
		var priority int
		parsedCount, err := fmt.Sscanf(args, "%s %s %d", &name, &url, &priority)
		if err != nil || parsedCount != 3 {
			return errors.New("некорректный формат. Используйте: \"Имя\" \"URL\" \"Приоритет\"")
		}

		source := model.Source{
			Name:     name,
			FeedURL:  url,
			Priority: priority,
		}

		sourceID, err := storage.Add(ctx, source)
		if err != nil {
			return err
		}

		reply := tgbotapi.NewMessage(update.Message.Chat.ID,
			fmt.Sprintf("Источник добавлен с ID: `%d`", sourceID))
		reply.ParseMode = parseModeMarkdownV2

		_, err = bot.Send(reply)
		return err
	}
}
