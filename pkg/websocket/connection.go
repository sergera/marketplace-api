package websocket

import (
	"context"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
	// time allowed to write a message to the peer
	writeWait = 10 * time.Second
)

type Connection struct {
	inner *websocket.Conn
	pool  *Pool
}

type Message struct {
	Type websocket.MessageType `json:"type"`
	Body string                `json:"body"`
}

func NewConnection(conn *websocket.Conn, pool *Pool) *Connection {
	return &Connection{conn, pool}
}

func (c *Connection) WriteText(body []byte) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(writeWait))
	defer cancel()
	err := c.inner.Write(ctx, websocket.MessageText, body)
	if err != nil {
		c.Close(websocket.StatusAbnormalClosure, "error writing text message")
		return err
	}
	return nil
}

func (c *Connection) WriteBinary(body []byte) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(writeWait))
	defer cancel()
	err := c.inner.Write(ctx, websocket.MessageBinary, body)
	if err != nil {
		c.Close(websocket.StatusAbnormalClosure, "error writing binary message")
		return err
	}
	return nil
}

func (c *Connection) WriteJSON(body interface{}) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(writeWait))
	defer cancel()
	err := wsjson.Write(ctx, c.inner, body)
	if err != nil {
		c.Close(websocket.StatusAbnormalClosure, "error writing JSON message")
		return err
	}
	return nil
}

func (c *Connection) Close(status websocket.StatusCode, reason string) {
	c.pool.Unregister <- c
	c.inner.Close(status, reason)
}
