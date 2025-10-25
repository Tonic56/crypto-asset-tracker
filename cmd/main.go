package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/WWoi/web-parcer/internal/aggregator"
	"github.com/WWoi/web-parcer/internal/models"
	"github.com/WWoi/web-parcer/internal/processor"
	"github.com/WWoi/web-parcer/internal/websocket"
)

func main() {
	fmt.Println("crypto-asset-tracker starting")

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	go func() {
		<-sigs
		cancel()
	}()

	// Channels
	rawMessages := make(chan []byte, 100)
	procOut := make(chan models.UniversalTrade, 100)

	// Websocket client (url can be changed)
	ws := websocket.New("wss://stream.binance.com:9443/stream?streams=btcusdt@aggTrade/ethusdt@aggTrade/bnbusdt@aggTrade", rawMessages, 5*time.Second)
	go ws.Start(ctx)

	// Processor (bytes -> models.UniversalTrade)
	proc := processor.New(rawMessages, procOut)
	go proc.Start(ctx)

	windowsChan := make(chan *models.Window)
	// Aggregator skeleton — wire input channel
	agg := aggregator.NewWindowAggregator(procOut, windowsChan)
	go agg.Start(ctx)

	go func() {
		for window := range windowsChan {
			fmt.Printf(
				"✅ Candle close: %s [%s] | Open: %.2f, High: %.2f, Low: %.2f, Close: %.2f | Volume: %f\n",
				window.Symbol,
				window.Interval,
				window.Open,
				window.High,
				window.Low,
				window.Close,
				window.Quantity,
			)
		}
	}()
	// Wait until context canceled
	<-ctx.Done()
	fmt.Println("\nshutting down")

	// задержка для graceful shutdown (опционально)
	time.Sleep(100 * time.Millisecond)
}
