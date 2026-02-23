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
	ID      int64  `json:"id" redis:"id"`
	ChatID  int64  `json:"chat_id" redis:"chat_id"`
	Balance int64  `json:"balance" redis:"balance"`
	Tariff  string `json:"tariff" redis:"tariff"`
	State   string `json:"state" redis:"state"`
}

var rdbpass = os.Getenv("REDIS_PASS")
var ctx = context.Background()

var acc_db = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	DB:       3,
	Password: rdbpass,
})

func main() {
	if rdbpass == "" {
		panic("REDIS_PASS environment variable not set")
	}
	http.HandleFunc("/speedball-api/v1/init", initHandler)
	http.HandleFunc("/speedball-api/v1/account", accountHandler)
	http.HandleFunc("/speedball-api/v1/keys", keysHandler)
	http.HandleFunc("/speedball-api/v1/topup", topUpHandler)
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
	account.ChatID, _ = strconv.ParseInt(r.URL.Query().Get("chat_id"), 10, 64)
	if account.ID == 0 || account.ChatID == 0 {
		http.Error(w, "Missing id, or chat_id parameter", http.StatusBadRequest)
		return
	}
	account.Balance = 0
	account.Tariff = "Стандартный"
	account.State = "Отключен"

	exist := acc_db.HExists(ctx, "user:"+strconv.FormatInt(account.ID, 10), "create_time").Val()
	if exist {
		w.WriteHeader(http.StatusAlreadyReported)
		return
	}
	acc_db.HSet(ctx, "user:"+strconv.FormatInt(account.ID, 10), &account)
	acc_db.HSet(ctx, "user:"+strconv.FormatInt(account.ID, 10), "create_time", time.Now().Unix())
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

func keysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	keys := acc_db.LRange(ctx, "user:"+id+":keys", 0, -1).Val()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func topUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID     string `json:"id"`
		Amount int64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.ID == "" || req.Amount <= 0 {
		http.Error(w, "Missing or invalid id or amount", http.StatusBadRequest)
		return
	}

	acc_db.HIncrBy(ctx, "user:"+req.ID, "balance", req.Amount)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok\n"))
}
