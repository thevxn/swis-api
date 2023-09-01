package config

import (
	"encoding/json"

	gin "github.com/gin-gonic/gin"
)

// JSONLogMiddleware transforms plaintext-formatted gin logs into JSON streams.
// https://stackoverflow.com/a/73936927
func JSONLogMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency

			s, err := json.Marshal(log)
			if err != nil {
				panic(err)
			}
			return string(s) + "\n"
		},
	)
}
