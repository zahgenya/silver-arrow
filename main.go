package main

import (
	"fmt"

	"silver-arrow/internal/streamer"
)

func main() {
	err := streamer.StartMiniTickerStream(streamer.Symbol, streamer.StreamDuration)
	if err != nil {
		fmt.Println("Application failed:", err)
	}
}
