package bot

import (
	"fmt"
	"log"
	"os"
	"time"
)

var quit chan struct{}

func startTimer(seconds int) {
	if quit != nil {
		close(quit)
	}

	quit = make(chan struct{})

	go func() {
		filename := "F:/Media/Twitch/Bot/timer.txt"
		file, err := os.Create(filename)

		if err != nil {
			log.Fatal(err)
		}

		defer file.Close()
		countdown := time.Duration(seconds) * time.Second
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		log.Printf("bot: timer started with duration %d seconds", seconds)

		for countdown >= 0 {
			select {
			case <-ticker.C:
				totalSeconds := int(countdown.Seconds())
				minutes := totalSeconds / 60
				seconds := totalSeconds % 60
				countdownMsg := fmt.Sprintf("%02d:%02d", minutes, seconds)

				file.Seek(0, 0)
				_, err = file.WriteString("")

				if err != nil {
					log.Fatal(err)
				}

				_, err = file.WriteString(countdownMsg)

				if err != nil {
					log.Fatal(err)
				}

				log.Printf("bot: timer updated to %s", countdownMsg)

				countdown -= time.Second
			case <-quit:
				log.Println("bot: timer stopped")

				return
			}
		}
	}()
}
