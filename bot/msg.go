package main

import (
	"fmt"
	"speedball/bot/accounts"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)


func HomeMsg(a accounts.Account, u string) *bot.SendMessageParams {
	msgText := "Бот управления доступом SurfBoost VPN" + "\n\n" +
		"Пользователь " + u + "!\n\n" +
		"Твой баланс: " + fmt.Sprintf("%d", a.Balance) + "\n" +
		"Тариф: " + a.Tariff + "\n" +
		"Статус доступа к VPN: " + a.State + "\n"

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{{Text: "⚙️ Подключение к VPN", CallbackData: "vpnConnect"}},
			{{Text: "💲 Внесение оплаты за VPN", CallbackData: "paymentMenu"}},
			{{Text: "💵 Акция «Приведи друга»", CallbackData: "referral"}},
			{{Text: "💸 Пожертвовать", CallbackData: "donate"}},
			{{Text: "💬 Помощь", CallbackData: "help"}},
		},
	}

	return &bot.SendMessageParams{
		ChatID:      a.ChatID,
		Text:        msgText,
		ReplyMarkup: keyboard,
	}
}
