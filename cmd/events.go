package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ChannelChatMsgSubPayload struct {
	Payload struct {
		Event struct {
			Msg struct {
				Text string `json:"text"`
			} `json:"message"`
		} `json:"event"`
	} `json:"payload"`
}

func ChannelChatMsgSub(authToken, session_id string) {
	config := readConfig()

	data := map[string]interface{}{
		"type":    "channel.chat.message",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": config.BroadcasterUserId,
			"user_id":             config.BroadcasterUserId,
		},
		"transport": map[string]string{
			"method":     "websocket",
			"session_id": session_id,
		},
	}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	eventSub(authToken, jsonData)
}

func eventSub(authToken string, subData []byte) {
	baseUrl := &url.URL{
		Scheme: "https",
		Host:   "api.twitch.tv",
		Path:   "helix/eventsub/subscriptions",
	}

	req, err := http.NewRequest("POST", baseUrl.String(), bytes.NewBuffer(subData))

	if err != nil {
		log.Fatal(err)
	}

	config := readConfig()

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Client-Id", config.ClientId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(body))
}
