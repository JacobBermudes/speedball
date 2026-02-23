package accounts

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Account struct {
	ID      int64  `json:"id"`
	ChatID  int64  `json:"chat_id"`
	Balance int64  `json:"balance"`
	Tariff  string `json:"tariff"`
	State   string `json:"state"`
}

func (a *Account) Init(r string) {
	initParams := url.Values{
		"id":      {strconv.FormatInt(a.ID, 10)},
		"chat_id": {strconv.FormatInt(a.ChatID, 10)},
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

	if r != "" && resp.StatusCode == http.StatusCreated {
		var req struct {
			ID     string `json:"id"`
			Amount int64 `json:"amount"`
		}
		req.ID = strconv.FormatInt(a.ID, 10)
		req.Amount = 55
		jsonData, err := json.Marshal(req)
		if err != nil {
			panic(err)
		}

		respA, err := http.Post("http://localhost:8801/speedball-api/v1/topup", "application/json", strings.NewReader(string(jsonData)))
		if err != nil {
			panic(err)
		}
		defer respA.Body.Close()
		req.ID = r
		jsonData, err = json.Marshal(req)
		if err != nil {
			panic(err)
		}
		respB, err := http.Post("http://localhost:8801/speedball-api/v1/topup", "application/json", strings.NewReader(string(jsonData)))
		if err != nil {
			panic(err)
		}
		defer respB.Body.Close()

		if resp.StatusCode != http.StatusOK {
			panic("Failed to send referral bonus")
		}
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
