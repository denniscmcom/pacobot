package bot

import (
	"log"
	"strconv"
)

func HandleCmd(cmd []string) {
	cmdReceived := cmd[0]
	log.Printf("bot: %s command received", cmdReceived)

	switch cmdReceived {
	case "timer":
		seconds, err := strconv.Atoi(cmd[1])

		if err != nil {
			log.Fatal("bot: invalid command arguments")
		}

		startTimer(seconds)
	}
}
