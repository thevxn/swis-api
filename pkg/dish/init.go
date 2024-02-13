package dish

import (
	"time"
)

var (
	EventChannel chan Message
)

func init() {
	go heartbeat()
}

func heartbeat() {
	for {
		if time.Now().Unix()%30 == 0 {
			//stream.Message <- Message{
			EventChannel <- Message{
				Content:    "heartbeat",
				SocketList: []string{},
				Timestamp:  time.Now().Unix(),
			}

			time.Sleep(time.Second * 1)
		}
	}
}
