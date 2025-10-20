package websocket

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type WSclient struct {
	url            string
	conn           *websocket.Conn
	outputChan     chan<- []byte
	stopChan       chan struct{}
	reconnectDelay time.Duration
}

func New(url string, output chan<- []byte, time time.Duration) *WSclient {
	return &WSclient{
		url:            url,
		outputChan:     output,
		reconnectDelay: time,
		stopChan:       make(chan struct{}),
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
			log.Println("WebSocket client stopped")
			return
		default:
			err := c.connect()
			if err != nil {
				log.Printf("Connection failed: %v. Retrying in %v", err, currentDelay)

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
			log.Println("✅ WebSocket connected successfully")

			go c.setupPingPong()
			c.readMessage(ctx) // Передаем контекст

			// Соединение разорвано, ждем перед переподключением
			select {
			case <-time.After(2 * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}
}

// Удаляем Reconnect или делаем ее приватной
// func (c *WSclient) reconnect(ctx context.Context) {

// }

func (c *WSclient) connect() error {
	conn, _, err := websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// readMessage тоже принимает context
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
			c.conn.SetReadDeadline(time.Now().Add(1 * time.Second))

			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					slog.Error("Read message error", slog.String("error", err.Error()))
				}
				return // Выходим при ошибке чтения
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
