package event

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"github.com/denniscmcom/pacobot/auth"
)

func ChannelChatMsgSub(authToken, session_id string) {
	config := auth.ReadConfig()

	data := map[string]any{
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

	log.Printf("event: subscribing to %s", data["type"])
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

	config := auth.ReadConfig()

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Client-Id", config.ClientId)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("status code: %d", res.StatusCode)

	if res.StatusCode != 202 {
		log.Fatal("event: failed to subscribe to event")
	}

	log.Println("event: subscribed")

	defer res.Body.Close()
}
