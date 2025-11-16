package api

import (
	"sync"
)

type CoinPrices struct {
	mu			sync.RWMutex
	lastPrice	float64
	baseVolume	float64
}

func NewCoinPrices() *CoinPrices {
	return &CoinPrices{}
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