package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"speedball/bot/commandHandlers"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDebug(),
		bot.WithAllowedUpdates([]string{"message", "callback_query"}),
		bot.WithWebhookSecretToken("token"),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		log.Fatal("Bot create FAIL:", err)
	}

	webhookAddr := "/speedball-webhook"
	webhookURL := "https://" + speedball_domen + ":443" + webhookAddr

	if _, err := b.SetWebhook(ctx, &bot.SetWebhookParams{
		URL:            webhookURL,
		AllowedUpdates: []string{"message", "callback_query"},
		SecretToken:    "token",
	}); err != nil {
		log.Fatal("Setting webhook FAIL:", err)
	}
	log.Println("Webhook setted:", webhookURL)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypePrefix,
		commandhandlers.StartHandler,
	)

	go b.StartWebhook(ctx)

	http.HandleFunc(webhookAddr, b.WebhookHandler())
	http.HandleFunc("/speedball-notify", notifyHandler(b))

	log.Println("Speedball WebHook listening :8800 (HTTP)")
	if err := http.ListenAndServe(":8800", nil); err != nil {
		log.Fatal("HTTP WebHook-Server FAULT:", err)
	}
}

// notifyHandler returns handler for internal POST notifications
func notifyHandler(b *bot.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		type internalSendReq struct {
			Chat_id string `json:"chat_id"`
			Text    string `json:"text"`
		}
		var req internalSendReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			fmt.Printf("BAD notify JSON")
			return
		}
		if req.Chat_id == "" || strings.TrimSpace(req.Text) == "" {
			fmt.Printf("missing cid/text")
			return
		}
		cid, _ := strconv.ParseInt(req.Chat_id, 10, 64)
		_, err := b.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: cid,
			Text:   req.Text,
		})
		if err != nil {
			log.Println("send fail:", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
