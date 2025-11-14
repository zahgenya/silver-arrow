package streamer

import (
	"fmt"
	"time"
	"strconv"

	binanceCon "github.com/binance/binance-connector-go"
	"silver-arrow/api"
)

const (
	Symbol         = "ETHUSDT"
	StreamDuration = 30 * time.Second
	// BASE_WS_URL default
)

func errHandler(err error) {
	fmt.Println("WebSocket Error:", err)
}

func StartMiniTickerStream(symbol string, duration time.Duration, coin *api.CoinPrices) error {
	wsStreamClient := binanceCon.NewWebsocketStreamClient(false)

	fmt.Printf("Subscribing to %s@miniTicker stream. Receiving updates for %s...\n", symbol, duration)

	wsMarketMiniTickerHandler := func(event binanceCon.WsMarketMiniTickerStatEvent) {
		fmt.Printf("Symbol: %s | Last Price (c): %s | Base Volume (v): %s\n",
			event.Symbol, event.LastPrice, event.BaseVolume)

		lastPrice, err1 := strconv.ParseFloat(event.LastPrice, 64)
        baseVol, err2 := strconv.ParseFloat(event.BaseVolume, 64)
        if err1 != nil || err2 != nil {
            fmt.Printf("parse error: %v %v\n", err1, err2)
            return
        }
		coin.Set(lastPrice, baseVol)
	}

	doneCh, stopCh, err := wsStreamClient.WsMarketMiniTickersStatServe(symbol, wsMarketMiniTickerHandler, errHandler)
	if err != nil {
		return fmt.Errorf("error starting Mini-Ticker stream: %w", err)
	}

	go func() {
		time.Sleep(duration)
		fmt.Printf("\nStopping WebSocket stream after %s...\n", duration)

		stopCh <- struct{}{}
	}()

	<-doneCh
	fmt.Println("Stream closed gracefully.")
	return nil
}
