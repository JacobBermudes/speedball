package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Account struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	ChatID   int64  `json:"chat_id"`
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
	account.ID, _ = strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	account.Username = r.URL.Query().Get("username")
	account.ChatID, _ = strconv.ParseInt(r.URL.Query().Get("chat_id"), 10, 64)
	if account.ID == 0 || account.Username == "" || account.ChatID == 0 {
		http.Error(w, "Missing id, username or chat_id parameter", http.StatusBadRequest)
		return
	}

	exist := acc_db.HExists(ctx, "user:"+strconv.FormatInt(account.ID, 10), "create_time").Val()
	if exist {
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}
	acc_db.HSet(ctx, "user:"+strconv.FormatInt(account.ID, 10), account, "create_time", time.Now().Format("02.01.2006"))
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
