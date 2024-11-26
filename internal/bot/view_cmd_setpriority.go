package bot

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"AI-RSS-Telegram-Bot/internal/botkit"
)

type PrioritySetter interface {
	SetPriority(ctx context.Context, sourceID int64, priority int) error
}

func ViewCmdSetPriority(prioritySetter PrioritySetter) botkit.ViewFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
		args := update.Message.CommandArguments()
		var sourceID int64
		var priority int
		parsedCount, err := fmt.Sscanf(args, "%d %d", &sourceID, &priority)
		if err != nil || parsedCount != 2 {
			return errors.New("некорректный формат. Используйте:  \"URL\" \"Приоритет\"")
		}

		if err := prioritySetter.SetPriority(ctx, sourceID, priority); err != nil {
			return err
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Приоритет успешно обновлен")

		if _, err := bot.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
