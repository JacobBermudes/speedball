package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	commandhandlers "speedball/bot/commandHandlers"
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

	<-ctx.Done()
	log.Println("Shutting down...")
}

// notifyHandler returns handler for internal POST notifications
func notifyHandler(b *bot.Bot) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[NOTIFY]: %s %s", r.Method, r.URL.Path)

		if r.Method != http.MethodPost {
			log.Println("[NOTIFY] Method not allowed")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			Chat_id string `json:"chat_id"`
			Text    string `json:"text"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("[NOTIFY] JSON decode error: %v", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		chatIDStr := strings.TrimSpace(req.Chat_id)
		text := strings.TrimSpace(req.Text)

		if chatIDStr == "" || text == "" {
			log.Println("[NOTIFY] Missing chat_id or text")
			http.Error(w, "chat_id and text are required", http.StatusBadRequest)
			return
		}

		chatID, err := strconv.ParseInt(chatIDStr, 10, 64)
		if err != nil {
			log.Printf("[NOTIFY] Invalid chat_id: %v", err)
			http.Error(w, "chat_id must be integer", http.StatusBadRequest)
			return
		}

		log.Printf("[NOTIFY] chat_id=%d, text=%q", chatID, text)

		_, sendErr := b.SendMessage(context.Background(), &bot.SendMessageParams{
			ChatID: chatID,
			Text:   text,
		})

		if sendErr != nil {
			log.Printf("[NOTIFY] SendMessage error: %v", sendErr)
			http.Error(w, fmt.Sprintf("Telegram error: %v", sendErr), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok\n"))
	}
}
