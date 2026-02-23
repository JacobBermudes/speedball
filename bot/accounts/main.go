package accounts

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type Account struct {
	ID       int64  `json:"id"`
	ChatID   int64  `json:"chat_id"`
	Balance  int64  `json:"balance"`
	Tariff   string `json:"tariff"`
	State    string `json:"state"`
}

func (a *Account) Init() {
	initParams := url.Values{
		"id":       {strconv.FormatInt(a.ID, 10)},
		"chat_id":  {strconv.FormatInt(a.ChatID, 10)},
	}
	initUrl := "http://localhost:8801/speedball-api/v1/init?"

	resp, err := http.Get(initUrl + initParams.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAlreadyReported && resp.StatusCode != http.StatusCreated {
		panic("Failed to initialize account")
	}
}

func (a *Account) GetData() {
	ap := url.Values{
		"id": {strconv.FormatInt(a.ID, 10)},
	}
	resp, err := http.Get("http://localhost:8801/speedball-api/v1/account?" + ap.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Failed to get account data")
	}

	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		panic(err)
	}
}

func (a *Account) GetKeys() []string {
	ap := url.Values{
		"id": {strconv.FormatInt(a.ID, 10)},
	}
	resp, err := http.Get("http://localhost:8801/speedball-api/v1/keys?" + ap.Encode())
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic("Failed to get account keys")
	}

	var keys []string
	if err := json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		panic(err)
	}
	return keys
}