package cmd

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type MetadataRes struct {
	Metadata struct {
		MsgType string `json:"message_type"`
		SubType string `json:"subscription_type"`
	} `json:"metadata"`
}

type WelcomeMsgPayload struct {
	Payload struct {
		Session struct {
			Id string `json:"id"`
		} `json:"session"`
	} `json:"payload"`
}

func ConnSocket(authToken string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	baseUrl := url.URL{Scheme: "wss", Host: "eventsub.wss.twitch.tv", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(baseUrl.String(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	done := make(chan struct{})

	ticker := time.NewTicker(time.Second * 15)
	defer ticker.Stop()

	go func() {
		defer close(done)

		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				log.Fatal(err)
			}

			var metadataRes MetadataRes

			if err := json.Unmarshal(msg, &metadataRes); err != nil {
				log.Fatal(err)
			}

			switch msgType := metadataRes.Metadata.MsgType; msgType {
			case "session_welcome":
				var welcomeMsgRes WelcomeMsgPayload

				if err := json.Unmarshal(msg, &welcomeMsgRes); err != nil {
					log.Fatal(err)
				}

				ChannelChatMsgSub(authToken, welcomeMsgRes.Payload.Session.Id)
			case "session_keepalive":
				ticker.Reset(time.Second * 15)
			case "notification":
				switch subType := metadataRes.Metadata.SubType; subType {
				case "channel.chat.message":
					var channelChatMsgSubPayload ChannelChatMsgSubPayload

					if err := json.Unmarshal(msg, &channelChatMsgSubPayload); err != nil {
						log.Fatal(err)
					}

					log.Println(string(msg))
					log.Println(channelChatMsgSubPayload.Payload.Event.Msg.Text)

				}
			default:
				log.Fatalf("%s: message type not implemented", msgType)
			}

		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			err := conn.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

			if err != nil {
				log.Fatal(err)
			}

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		case <-ticker.C:
			// TODO: Replace this with logic to reconnect
			log.Fatal("connection closed: timeout")
		}
	}
}
