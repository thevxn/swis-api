package main

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
)

// https://stackoverflow.com/a/29439630
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://swjango.savla.su")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Auth-Token")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, UPDATE")

		if c.Request.Method == "OPTIONS" {
			//c.AbortWithStatus(http.StatusOK)
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
