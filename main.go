package main

import (
	"fmt"
	"net/http"

	"silver-arrow/internal/streamer"
	"silver-arrow/api"
	"github.com/gin-gonic/gin"
)

func main() {
	coin := &api.CoinPrices{}

	go func() {
		err := streamer.StartMiniTickerStream(streamer.Symbol, streamer.StreamDuration, coin)
		if err != nil {
			fmt.Println()
		}
	}()

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/latest-price", api.GetLatestPrice(coin))
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	fmt.Println("Server listening on http://localhost:8080")
	router.Run(":8080")
}
