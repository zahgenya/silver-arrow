package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	Prices	*CoinPrices
}

func NewApiHandler(prices *CoinPrices) *ApiHandler {
	return &ApiHandler{
		Prices: prices,
	}
}

func (a *ApiHandler) GetLatestPrice(c *gin.Context) {
	latestPrice, baseVolume := a.Prices.Get()

	if latestPrice == 0.0 && baseVolume == 0.0 {
		response := PriceNotAvailableResponse{
			Message: "no price avaliable yet",
			Symbol:  "ETHUSDT",
		}

		c.JSON(http.StatusAccepted, response)
		return
	}

	response := LatestPriceResponse{
		Symbol:      "ETHUSDT",
		LatestPrice: latestPrice,
		BaseVolume:  baseVolume,
	}
	c.JSON(http.StatusOK, response)
}

func (a *ApiHandler) GetHealth(c *gin.Context) {
	response := HealthResponse{
		Status: "ok",
	}
	c.JSON(http.StatusOK, response)
}