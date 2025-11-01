package websocket

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	AggTradeURL = "wss://stream.binance.com:9443/stream?streams=btcusdt@aggTrade/ethusdt@aggTrade/bnbusdt@aggTrade"

	// All - все монеты, Several - определенные
	// @3000 -> присылает окно каждые 3 секунды (хотя по факту куда реже)
	MiniTickerAllURL     = "wss://stream.binance.com:443/ws/!miniTicker@arr@1000ms"
	MiniTickerSeveralURL = "wss://stream.binance.com:9443/stream?streams=btcusdt@miniTicker/ethusdt@miniTicker/bnbusdt@miniTicker"
)

const (
	AggTrade   = "aggTrade"
	MiniTicker = "24hrMiniTicker"
)

type WSclient struct {
	url            string
	conn           *websocket.Conn
	outputChan     chan<- []byte
	reconnectDelay time.Duration
}

func New(url string, output chan<- []byte, reconnectDelay time.Duration) *WSclient {
	return &WSclient{
		url:            url,
		outputChan:     output,
		reconnectDelay: reconnectDelay,
	}
}

// принимаем context
func (c *WSclient) Start(ctx context.Context) {
	currentDelay := 1 * time.Second
	maxDelay := 2 * time.Minute

	for {
		select {
		case <-ctx.Done():
			if c.conn != nil {
				c.conn.Close()
			}
			slog.Info("WebSocket client stopped")
			return
		default:
			err := c.connect()
			if err != nil {
				slog.Error(
					"❌ Connection failed",
					"error", err,
					"retry_in", currentDelay,
					"url", c.url,
				)

				// Ждем с проверкой контекста
				select {
				case <-time.After(currentDelay):
					// Продолжаем
				case <-ctx.Done():
					return
				}

				// Экспоненциальная задержка
				currentDelay *= 2
				if currentDelay > maxDelay {
					currentDelay = maxDelay
				}
				continue
			}

			// Успешное подключение
			currentDelay = 1 * time.Second
			slog.Info("✅ WebSocket connected successfully")

			go c.setupPingPong(ctx)
			c.readMessage(ctx)

			// Соединение разорвано, ждем перед переподключением

			select {
			case <-time.After(2 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}
}

func (c *WSclient) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *WSclient) readMessage(ctx context.Context) {
	defer func() {
		if c.conn != nil {
			c.conn.Close()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping message reader")
			return
		default:
			// Устанавливаем таймаут для чтения, чтобы можно было проверить контекст
			c.conn.SetReadDeadline(time.Now().Add(5 * time.Minute))

			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(
					err,
					websocket.CloseGoingAway,
					websocket.CloseAbnormalClosure,
				) {
					slog.Error("❌ Read message error", slog.String("error", err.Error()))
				}
				return
			}

			// Отправляем сообщение в канал
			select {
			case c.outputChan <- msg:
			case <-ctx.Done():
				return
			}
		}
	}
}
