package streamer

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"silver-arrow/api"

	binanceCon "github.com/binance/binance-connector-go"
)

type MiniTickerEvent binanceCon.WsMarketMiniTickerStatEvent

type Stream struct {
	doneCh <-chan struct{}
	stopCh chan<- struct{}
}

func (s *Stream) Close() {
	fmt.Println("Closing stream...")
	s.stopCh <- struct{}{}
	<-s.doneCh
	fmt.Println("Stream closed.")
}

func errHandler(err error) {
	fmt.Println("WebSocket Error:", err)
}

func NewMiniTickerStream(symbol string) (<-chan MiniTickerEvent, *Stream, error) {
	wsStreamClient := binanceCon.NewWebsocketStreamClient(false)
	eventCh := make(chan MiniTickerEvent)

	wsMarketMiniTickerHandler := func(event binanceCon.WsMarketMiniTickerStatEvent) {
		eventCh <- MiniTickerEvent(event)
	}

	doneCh, stopCh, err := wsStreamClient.WsMarketMiniTickersStatServe(symbol, wsMarketMiniTickerHandler, errHandler)
	if err != nil {
		close(eventCh)
		return nil, nil, fmt.Errorf("failed to start WebSocket stream: %v", err)
	}

	stream := &Stream{
		doneCh: doneCh,
		stopCh: stopCh,
	}

	return eventCh, stream, nil
}

func ProcessMiniTickerEvents(ctx context.Context, symbol string, eventCh <-chan MiniTickerEvent, coin *api.CoinPrices) {
	for {
		select {
		case event := <-eventCh:
			lastPrice, err := strconv.ParseFloat(event.LastPrice, 64)
			if err != nil {
				log.Printf("ERROR: [%s] Price parsing failed: %v. Raw value: %q", event.Symbol, err, event.LastPrice)
				continue
			}

			baseVol, err := strconv.ParseFloat(event.BaseVolume, 64)
			if err != nil {
				log.Printf("ERROR: [%s] Volume parsing failed: %v. Raw value: %q", event.Symbol, err, event.BaseVolume)
				continue
			}

			coin.Set(lastPrice, baseVol)

		case <-ctx.Done():
			return
		}
	}
}
