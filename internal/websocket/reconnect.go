package websocket

// import (
// 	"log"
// 	"time"
// )


// func (c *WSclient) Reconnect() {
// 	currentDelay := 1 * time.Second
// 	maxDelay := 2 * time.Minute

// 	for {
// 		err := c.connect()

// 		if err != nil {
// 			log.Printf("Error: %v. Retry through %v", err, currentDelay)
// 			time.Sleep(currentDelay)

// 			// Увеличиваем задержку в 2 раза 
// 			currentDelay *= 2
// 			if currentDelay > maxDelay {
// 				currentDelay = maxDelay
// 			}

// 			continue
// 		}

// 		// Успешно подключились - сбрасываем задержку
// 		currentDelay = 1 * time.Second
// 		log.Println("✅ Connection succsessful")

// 		c.setupPingPong()
// 		c.readMessage()
// 		c.conn.Close()

// 		time.Sleep(2 * time.Second)
// 	}
// }
