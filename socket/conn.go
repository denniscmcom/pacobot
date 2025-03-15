package socket

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/denniscmcom/pacobot/bot"
	"github.com/denniscmcom/pacobot/event"
	"github.com/gorilla/websocket"
)

func Connect(authToken string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	baseUrl := url.URL{Scheme: "wss", Host: "eventsub.wss.twitch.tv", Path: "/ws"}

	log.Println("socket: connecting...")
	conn, _, err := websocket.DefaultDialer.Dial(baseUrl.String(), nil)

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	log.Println("socket: connected")

	var timeout time.Ticker
	done := make(chan struct{})

	go readMsg(done, conn, &timeout, authToken)

	for {
		select {
		case <-interrupt:
			closeConn(conn)

			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		case <-done:
			log.Println("socket: connection closed by server")
			Connect(authToken)
		case <-timeout.C:
			log.Println("socket: connection lost")
			timeout.Stop()
			Connect(authToken)
		}
	}
}

func readMsg(done chan struct{}, conn *websocket.Conn, timeout *time.Ticker, authToken string) {
	defer close(done)
	var timeout_secs time.Duration

	for {
		log.Println("socket: waiting for msg...")
		_, msg, err := conn.ReadMessage()

		if err != nil {
			break
		}

		var resMetadata Res_Metadata

		if err := json.Unmarshal(msg, &resMetadata); err != nil {
			log.Fatal(err)
		}

		msgType := resMetadata.Metadata.MsgType
		log.Printf("socket: %s msg received", msgType)

		switch msgType {
		case "session_welcome":
			var resWelcome Res_Welcome

			if err := json.Unmarshal(msg, &resWelcome); err != nil {
				log.Fatal(err)
			}

			timeout_secs = time.Duration(resWelcome.Payload.Session.KeepaliveTimeout+3) * time.Second
			timeout = time.NewTicker(timeout_secs)
			defer timeout.Stop()

			event.ChannelChatMsgSub(authToken, resWelcome.Payload.Session.Id)

		case "session_keepalive":
			timeout.Reset(timeout_secs)
			log.Println("socket: timeout resetted")

		case "notification":
			var resMetadataNotif Res_Metadata_Notif

			if err := json.Unmarshal(msg, &resMetadataNotif); err != nil {
				log.Fatal(err)
			}

			subType := resMetadataNotif.Metadata.SubType
			log.Printf("socket: %s event received", subType)

			switch subType {
			case "channel.chat.message":
				var resNotifChannelChatMsg Res_Notif_ChannelChatMsg

				if err := json.Unmarshal(msg, &resNotifChannelChatMsg); err != nil {
					log.Fatal(err)
				}

				chatMsg := resNotifChannelChatMsg.Payload.Event.Msg.Text

				if strings.HasPrefix(chatMsg, "!") {
					go bot.HandleCmd(strings.Split(chatMsg[1:], " "))
				}
			}
		default:
			log.Fatalf("socket: %s message type not implemented", msgType)
		}
	}
}

func closeConn(conn *websocket.Conn) {
	err := conn.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("socket: connection closed")
}

// func test() {
// 	var res Response

// 	// Deserializas
// 	err := json.Unmarshal([]byte(jsonData), &res)

// 	if err != nil {
// 		fmt.Println("Error al deserializar:", err)
// 		return
// 	}

// 	// Conviertes la estructura nuevamente a JSON formateado

// 	formattedJSON, err := json.MarshalIndent(res, "", " ")

// 	if err != nil {
// 		fmt.Println("Error al formatear JSON:", err)
// 		return
// 	}
// }
