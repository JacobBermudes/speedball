package commandhandlers

import (
	"context"
	"speedball/bot/accounts"
	"speedball/bot/messages"
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

	var account accounts.Account
	account.ChatID = update.Message.Chat.ID
	account.ID = update.Message.From.ID

	var ref string
	if len(args) > 0 {
		ref = args[0]
	}

	account.Init(ref)

	b.SendMessage(ctx, messages.HomeMsg(account, update.Message.From.FirstName))
}
