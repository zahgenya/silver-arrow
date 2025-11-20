package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"silver-arrow/api"
	"silver-arrow/internal/streamer"
	"silver-arrow/internal/telegram"
)

//go:generate go tool oapi-codegen --config=codegen.yaml openapi.yaml

const (
	StreamDuration = 30 * time.Second
	APIPort        = ":8080"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
		return
	}
	//later change it to config struct pointer or something like that
	apiKey := os.Getenv("TELEGRAM_URL")
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

		router.GET("/openapi.yaml", func(c *gin.Context) {
			c.File("./openapi.yaml")
		})

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/openapi.yaml"),
		))

		ethPrices := coinPricesMap["ETHUSDT"]
		apiHandler := api.NewApiHandler(ethPrices)

		api.RegisterHandlers(router, apiHandler)

		fmt.Println("Server listening on http://localhost" + APIPort)
		fmt.Println("Swagger UI available at http://localhost" + APIPort + "/swagger/index.html")

		if err := router.Run(APIPort); err != nil && err != http.ErrServerClosed {
			fmt.Printf("API server error: %v\n", err)
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := telegram.StartTelegramPoll(apiKey)
		if err != nil {
			fmt.Printf("Error starting poll for telegram: %v\n", err)
			return
		}
	}()

	wg.Wait()
	fmt.Println("All streams shut down gracefully.")
}
