package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"silver-arrow/api"
	"silver-arrow/internal/streamer"

	"github.com/gin-gonic/gin"
)

const (
	StreamDuration = 30 * time.Second
	APIPort        = ":8080"
)

func main() {
	symbols := []string{"ETHUSDT", "SOLUSDT"}
	var wg sync.WaitGroup

	coinPricesMap := make(map[string]*api.CoinPrices)
	for _, sym := range symbols {
		coinPricesMap[sym] = &api.CoinPrices{}
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("Shutting down gracefully...")
		cancel()
	}()

	go func() {
		time.Sleep(StreamDuration)
		fmt.Printf("\nStopping streams after %s...\n", StreamDuration)
	}()

	for _, symbol := range symbols {
		wg.Add(1)
		go func(sym string, coin *api.CoinPrices) {
			defer wg.Done()

			eventCh, stream, err := streamer.NewMiniTickerStream(sym)
			if err != nil {
				fmt.Printf("Failed to start stream for %s: %v\n", sym, err)
				return
			}
			defer stream.Close()

			fmt.Printf("Subscribed to %s@miniTicker stream.\n", sym)

			streamer.ProcessMiniTickerEvents(ctx, sym, eventCh, coin)
		}(symbol, coinPricesMap[symbol])
	}

	go func() {
		router := gin.Default()

		apiV1 := router.Group("/api/v1")
		{
			ethCoin := coinPricesMap["ETHUSDT"]
			apiV1.GET("/latest-price", api.GetLatestPrice(ethCoin))
		}

		router.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		fmt.Println("Server listening on http://localhost" + APIPort)

		if err := router.Run(APIPort); err != nil {
			fmt.Printf("API server error: %v\n", err)
			cancel()
		}
	}()

	wg.Wait()
	fmt.Println("All streams shut down gracefully.")
}
