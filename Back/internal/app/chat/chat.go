package chat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan Message
	Hub  *Hub
}

func NewUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{ID: id, Conn: conn, Send: make(chan Message, 256), Hub: hub}
}

func (c *Client) Read(maxMessageSize int64, pongWait time.Duration) {

	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		c.Hub.broadcast <- msg
	}
}

func (c *Client) Write(pingPeriod time.Duration, writeWait time.Duration) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := c.Conn.WriteJSON(message)
				if err != nil {
					fmt.Printf("Error: %v \n", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) Close() {
	close(c.Send)
}

func ServeWS(ctx *gin.Context, roomId string, hub *Hub) error {
	const (
		writeWait      = 10 * time.Second
		pongWait       = 60 * time.Second
		pingPeriod     = (pongWait * 9) / 10
		maxMessageSize = 512
	)
	fmt.Print(roomId)

	up := NewUpgrader()
	ws, err := up.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return err
	}

	client := NewClient(roomId, ws, hub)

	hub.register <- client
	go client.Write(pingPeriod, writeWait)
	go client.Read(maxMessageSize, pongWait)

	return nil
}
