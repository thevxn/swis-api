package config

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	gin "github.com/gin-gonic/gin"
)

func MirrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// do not replicate GET requests
		if ctx.Request.Method == "GET" || ctx.Request.Method == "OPTIONS" {
			ctx.Next()
			return
		}

		// do not replicate already replicated API call
		if ctx.Request.Header.Get("X-Mirror-Request") != "" {
			/*ctx.AbortWithStatusJSON(http.StatusNoContent, gin.H{
				"code":    http.StatusNoContent,
				"message": "dropping already replicated API call",
			})*/
			ctx.Next()
			return
		}

		// we need to buffer the body if we want to read it here and send it
		// in the request.
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": fmt.Sprintf("cannot replicate request body: %s", err.Error()),
			})
			return
		}

		// you can reassign the body if you need to parse it as multipart
		ctx.Request.Body = io.NopCloser(bytes.NewReader(body))

		// create a new url from the raw RequestURI sent by the client
		proxyScheme := "http"
		proxyHost := "localhost:8051"
		url := fmt.Sprintf("%s://%s%s", proxyScheme, proxyHost, ctx.Request.RequestURI)

		newReq, err := http.NewRequest(ctx.Request.Method, url, bytes.NewReader(body))

		// We may want to filter some headers, otherwise we could just use a shallow copy
		newReq.Header = ctx.Request.Header
		/*newReq.Header = make(http.Header)
		for h, val := range ctx.Request.Header {
			newReq.Header.Set(h, val)
		}*/
		newReq.Header.Set("X-Mirror-Request", os.Getenv("HOSTNAME"))

		client := http.Client{}
		resp, err := client.Do(newReq)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadGateway, gin.H{
				"code":    http.StatusBadGateway,
				"message": fmt.Sprintf("cannot reach linked instance: %s", err.Error()),
			})
			return
		}
		defer resp.Body.Close()

		ctx.Next()
	}
}
