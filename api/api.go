package api

import (
	"sync"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CoinPrices struct {
	mu			sync.RWMutex
	lastPrice	float64
	baseVolume	float64
}

func (coin *CoinPrices) Set(price float64, volume float64) {
	coin.mu.Lock()
	defer coin.mu.Unlock()
	coin.lastPrice = price
	coin.baseVolume = volume
}

func (coin *CoinPrices) Get() (float64, float64) {
	coin.mu.RLock()
	defer coin.mu.RUnlock()
	lastPrice := coin.lastPrice
	baseVolume := coin.baseVolume
	return lastPrice, baseVolume
}

func GetLatestPrice(coin *CoinPrices) gin.HandlerFunc {
	return func(c *gin.Context) {
		latestPrice, baseVolume := coin.Get()

		if latestPrice == 0.0 && baseVolume == 0.0{
			c.JSON(http.StatusOK, gin.H{
				"message": 		"no price avaliable yet",
				"symbol":		"ETHUSDT",
				"latest-price":	gin.H{},
				"base-volume":	gin.H{},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"symbol":		"ETHUSDT",
			"latest-price": latestPrice,
			"base-volume":	baseVolume,
		})
	}
}
