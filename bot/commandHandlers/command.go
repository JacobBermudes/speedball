package commandhandlers

import (
	"context"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func StartHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	text := update.Message.Text
	args := strings.Fields(strings.TrimPrefix(text, "/start"))

	var ref string
	if len(args) > 0 {
		ref = args[0]
	}

	msg := "Привет!"
	if ref != "" {
		msg += "\nРеферал от: " + ref + " — бонус начислен!"
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   msg,
	})
}
