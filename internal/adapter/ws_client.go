package adapter

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	URL        string
	Conn       *websocket.Conn
	Dialer     *websocket.Dialer
	writeMutex sync.Mutex
}

func NewWSClient(url string) *WSClient {
	return &WSClient{
		URL:    url,
		Dialer: &websocket.Dialer{HandshakeTimeout: 5 * time.Second},
	}
}

func (c *WSClient) Connect(ctx context.Context, header http.Header) error {
	conn, _, err := c.Dialer.DialContext(ctx, c.URL, header)
	if err != nil {
		return err
	}
	c.Conn = conn
	return nil
}

func (c *WSClient) SendMessage(messageType int, data []byte) error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	return c.Conn.WriteMessage(messageType, data)
}

func (c *WSClient) ReceiveMessages(ctx context.Context, handler func(msgType int, msg []byte)) {
	go func() { // runs a single go routine to listen
		for {
			select {
			case <-ctx.Done():
				return
			default:
				msgType, msg, err := c.Conn.ReadMessage()
				if err != nil {
					log.Printf("ws read error: %v\n", err)
					time.Sleep(time.Second) // Should change later
					continue
				}
				handler(msgType, msg)
			}
		}
	}()
}

func (c *WSClient) Close() error {
	c.writeMutex.Lock()
	defer c.writeMutex.Unlock()
	if c.Conn != nil {
		return c.Conn.Close()
	}
	return nil
}
