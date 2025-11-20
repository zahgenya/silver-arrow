package telegram

import (
	"net/http"
	"fmt"
	"io"
	"encoding/json"
	"bytes"
	"log"
	"time"
)

const (
	PollTimeout = 30
)

func initBot(apiKey string) (TgUser, error) {
	url := apiKey + "getMe"
	log.Printf("URL FOR INIT BOT IS: %v\n", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return TgUser{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return TgUser{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return TgUser{}, fmt.Errorf("telegram api error: getMe was failed with status code %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return TgUser{}, err
	}

	var user TgUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return TgUser{}, err
	}

	return user, nil
}

func fetchLatestPrice() (LatestPriceResponse, error) {
	url := "http://localhost:8080/api/v1/latest-price"

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LatestPriceResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return LatestPriceResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return LatestPriceResponse{}, fmt.Errorf("internal api error: was failed with status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LatestPriceResponse{}, err
	}

	var price LatestPriceResponse
	err = json.Unmarshal(body, &price)
	if err != nil {
		return LatestPriceResponse{}, err
	}

	return price, nil
}

func sendMessage(chatId int64, text string, apiKey string) error {
	url := apiKey + "sendMessage"

	message := TgSendMessage{
		ChatId: chatId,
		Text: text,
	}

	reqBody, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("telegram api error: sendMessage was failed with status code: %v body: %v", resp.StatusCode, respBody)
	}

	return nil
}

func getUpdates(offset int, apiKey string) ([]TgUpdate, error) {
	url := fmt.Sprintf("%vgetUpdates?offset=%v&timeout=%v", apiKey, offset, PollTimeout)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("telegram api error: getUpdates was failed with status: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var updates TgGetUpdate
	err = json.Unmarshal(body, &updates)
	if err != nil {
		return nil, err
	}

	return updates.Result, nil
}

func StartTelegramPoll(url string) error {
	_, err := initBot(url)
	if err != nil {
		return err
	}
	nextOffset := 0

	for {
		updates, err := getUpdates(nextOffset, url)
		if err != nil {
			log.Printf("Error during polling telegram updates: retrying for 3 seconds...")
			time.Sleep(3 * time.Second)
			continue
		}

		for _, update := range updates {
			price, err := fetchLatestPrice()
			if err != nil {
				log.Printf("%v", err)
				continue
			}

			msg := fmt.Sprintf("Coin %v price changed\nBase volume: %v\nPrice: %v\n", price.Symbol, price.BaseVolume, price.LatestPrice)
			err = sendMessage(update.Message.Chat.Id, msg, url)
			if err != nil {
				log.Printf("%v", err)
				continue
			}

			if update.UpdateId >= nextOffset {
				nextOffset = update.UpdateId + 1
			}
		}
	}
}