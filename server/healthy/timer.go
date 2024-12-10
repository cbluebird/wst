package healthy

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	interval = os.Getenv("AUTO_SHUTDOWN_INTERVAL")
	target   = os.Getenv("TARGET")
	jwtToken = os.Getenv("JWT_TOKEN")
)

var ActiveNum int32 = 0

func Monitor() {
	shutdownDuration, _ := time.ParseDuration(interval)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	zeroDuration := 0 * time.Minute
	for {
		select {
		case <-ticker.C:
			if atomic.LoadInt32(&ActiveNum) == 0 {
				log.Println("No active connection")
				zeroDuration += 1 * time.Minute
				if zeroDuration >= shutdownDuration {
					sendShutdownRequest()
					return
				}
			} else {
				log.Println("Active connection number:", atomic.LoadInt32(&ActiveNum))
				zeroDuration = 0
			}
		}
	}
}

func sendShutdownRequest() {
	url := target + "/opsrequest"
	data := map[string]string{
		"operation": "shutdown",
		"jwt_token": jwtToken,
	}
	jsonData, _ := json.Marshal(data)
	_, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// Handle error
	}
}
