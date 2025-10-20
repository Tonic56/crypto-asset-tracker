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
	ws := websocket.New("wss://example.com/ws", rawMessages, 5*time.Second)
	go ws.Start(ctx)

	// Processor (bytes -> models.UniversalTrade)
	proc := processor.New(rawMessages, procOut)
	go proc.Start()

	// Aggregator skeleton — wire input channel
	agg := aggregator.New(procOut)
	go agg.Start()

	// Wait until context canceled
	<-ctx.Done()
	fmt.Println("shutting down")
	
	// задержка для graceful shutdown (опционально) 
	time.Sleep(100 * time.Millisecond)
}
