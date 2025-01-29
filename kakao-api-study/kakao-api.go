package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"msa-study/configs"
	"net/http"
)

type KakaoMessage struct {
	Receiver string `json:"receiver"`
	Text     string `json:"text"`
}

func sendKakaoMessage(receiver, text string) error {
	kakaoConfig := configs.LoadKakaoConfig()
	url := "https://kapi.kakao.com/v2/api/talk/memo/default/send"

	message := KakaoMessage{
		Receiver: receiver,
		Text:     text,
	}

	messageData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(messageData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+kakaoConfig.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status: %s, response: %s", resp.Status, string(body))
	}

	log.Printf("Message sent successfully: %s", string(body))
	return nil
}

func main() {
	receiver := "receiver_id" // Replace with actual receiver ID
	text := "Hello from Kakao API!"

	err := sendKakaoMessage(receiver, text)
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}
}
