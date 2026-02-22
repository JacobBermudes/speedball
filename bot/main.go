package speedball_tg_bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"speedball/bot/accounts"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	speedball_domen := os.Getenv("SPEEDBALL_DOMEN")
	if speedball_domen == "" {
		log.Fatal("SPEEDBALL_DOMEN environment variable not set")
	}
	token := os.Getenv("SPEEDBALL_TG_BOT_TOKEN")
	if token == "" {
		log.Fatal("SPEEDBALL_TG_BOT_TOKEN environment variable not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Bot create FAIL:", err)
	}
	bot.Debug = true
	log.Printf("Auth as: @%s", bot.Self.UserName)
	webhookAddr := "/speedball_webhook"
	webhookURL := "https://" + speedball_domen + ":8443" + webhookAddr
	webhook, _ := tgbotapi.NewWebhook(webhookURL)
	webhook.AllowedUpdates = []string{"message", "callback_query"}

	_, err = bot.Request(webhook)
	if err != nil {
		log.Fatal("Setting webhook FAIL:", err)
	}
	log.Println("Webhook setted:", webhookURL)
	updates := bot.ListenForWebhook(webhookAddr)

	go func() {
		http.HandleFunc("/speedball_notify", func(w http.ResponseWriter, r *http.Request) {
			type internalSendReq struct {
				Cid  string `json:"cid"`
				Text string `json:"text"`
			}
			var req internalSendReq
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				fmt.Printf("BAD notify JSON")
				return
			}
			if req.Cid == "" || strings.TrimSpace(req.Text) == "" {
				fmt.Printf("missing cid/text")
				return
			}
			cid, _ := strconv.ParseInt(req.Cid, 10, 64)
			msg := tgbotapi.NewMessage(cid, req.Text)
			if _, err := bot.Send(msg); err != nil {
				log.Println("send fail:", err)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})

		if err := http.ListenAndServe(":8800", nil); err != nil {
			log.Fatal("HTTP WebHook-Server FAULT:", err)
		}
		log.Println("Speedball WebHook listening :8800 (HTTP)")
	}()

	for update := range updates {
		log.Printf("Get update: %+v", update)

		if update.Message != nil && update.Message.IsCommand() {
			if update.Message.Command() == "start" {

				account := accounts.Account{
					Username: update.Message.From.UserName,
					ID:       update.Message.From.ID,
					ChatID:   update.Message.Chat.ID,
				}
				account.Init()
				account.GetData()

				bot.Send(HomeMsg(account))
			}
		}

		if update.Message != nil {

		}

		if update.CallbackQuery != nil {

			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
			bot.Request(callback)
		}
	}
}
