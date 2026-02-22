package speedball_tg_bot

import (
	"fmt"
	"speedball/bot/accounts"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HomeMsg(a accounts.Account) tgbotapi.MessageConfig {

	msgText := "Бот управления доступом SurfBoost VPN" + "\n\n" +
		"Пользователь " + a.Username + "!\n\n" +
		"Твой баланс: " + fmt.Sprintf("%d", a.Balance) + "\n" +
		"Тариф: " + a.Tariff + "\n" +
		"Статус доступа к VPN: " + a.State + "\n"
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Подключение к VPN", "vpnConnect"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💲 Внесение оплаты за VPN", "paymentMenu"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💵 Акция «Приведи друга»", "referral"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💸 Пожертвовать", "donate"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💬 Помощь", "help"),
		),
	)
	msg := tgbotapi.NewMessage(a.ChatID, msgText)
	msg.ReplyMarkup = keyboard
	return msg
}
