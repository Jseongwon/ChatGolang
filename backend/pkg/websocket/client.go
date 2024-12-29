
import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	Id         string
	Connection *websocket.Conn
	Pool       *Pool
	Mutex      *sync.Mutex
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Connection.Close()
	}()

	for {
		messageType, readMessage, err := c.Connection.ReadMessage()
		if err != nil {
			c.Pool.Unregister <- c
			c.Connection.Close()
			break
		}

		message := Message{Type: messageType, Body: string(readMessage), Sender: c}
		c.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}