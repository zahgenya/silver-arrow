package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"silver-arrow/api"
	"silver-arrow/internal/streamer"
)

//go:generate go tool oapi-codegen --config=codegen.yaml openapi.yaml

func main() {
	coin := api.NewCoinPrices()

	go func() {
		for {
			err := streamer.StartMiniTickerStream(streamer.Symbol, streamer.StreamDuration, coin)
			if err != nil {
				log.Printf("Streamer exited with error: %v. Reconnecting...", err)
			}
		}
	}()

	apiHandler := api.NewApiHandler(coin)

	router := gin.Default()

	router.GET("/openapi.yaml", func(c *gin.Context) {
		c.File("./openapi.yaml")
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/openapi.yaml"),
	))

	api.RegisterHandlers(router, apiHandler)

	fmt.Println("Server listening on http://localhost:8080")
	log.Fatal(router.Run(":8080"))
}