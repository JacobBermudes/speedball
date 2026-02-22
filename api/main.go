package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	ChatID   string `json:"chat_id"`
	Balance  int64  `json:"balance"`
	Tariff   string `json:"tariff"`
	State    string `json:"state"`
}

var REDIS_PASS = os.Getenv("REDIS_PASS")
var ctx = context.Background()

var acc_db = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	DB:       3,
	Password: REDIS_PASS,
})

func main() {
	http.HandleFunc("/speedball-api/v1/init", initHandler)
	http.HandleFunc("/speedball-api/v1/account", accountHandler)
	http.ListenAndServe(":8801", nil)
	fmt.Println("Server starting on :8801")
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	var account Account
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	account.ID = r.URL.Query().Get("id")
	account.Username = r.URL.Query().Get("username")
	account.ChatID = r.URL.Query().Get("chat_id")
	if account.ID == "" || account.Username == "" || account.ChatID == "" {
		http.Error(w, "Missing id, username or chat_id parameter", http.StatusBadRequest)
		return
	}

	exist := acc_db.HExists(ctx, "user:"+account.ID, "create_time").Val()
	if exist {
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}
	acc_db.HSet(ctx, "user:"+account.ID, account, "create_time", time.Now().Format("02.01.2006"))
	w.WriteHeader(http.StatusCreated)
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	var account Account
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	err := acc_db.HGetAll(ctx, "user:"+id).Scan(&account)
	if err == redis.Nil {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}
