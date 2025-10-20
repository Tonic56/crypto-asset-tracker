package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func (c *WSclient) setupPingPong() {
	c.conn.SetPingHandler(func(appData string) error {
		log.Println("Ping from Binance, answer pong")

		// Отправляем pong обратно
		err := c.conn.WriteControl(
			websocket.PongMessage,         // Тип: PONG
			[]byte(appData),               // Тот же payload
			time.Now().Add(1*time.Second), // Deadline
		)

		return err
	})
}
