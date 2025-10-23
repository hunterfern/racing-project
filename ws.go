package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				_ = c.conn.Close()
			}
		case msg := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					delete(h.clients, c)
					close(c.send)
					_ = c.conn.Close()
				}
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Hub) WSHandler(commands chan<- Command) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		c := &Client{conn: conn, send: make(chan []byte, 256)}
		h.register <- c

		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer func() {
				ticker.Stop()
				h.unregister <- c
			}()
			for {
				select {
				case msg, ok := <-c.send:
					if !ok {
						_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
						return
					}
					if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
						return
					}
				case <-ticker.C:
					_ = c.conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second))
				}
			}
		}()

		for {
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				break
			}
			var m map[string]any
			if json.Unmarshal(data, &m) == nil {
				t, _ := m["type"].(string)
				kind := strings.ToUpper(t)
				args := map[string]string{}

				for k, v := range m {
					if k == "type" {
						continue
					}
					if s, ok := v.(string); ok {
						args[k] = s
					}

					if f, ok := v.(float64); ok {
						args[k] = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.0f", f), "0"), ".")
					}
				}
				if kind != "" {
					commands <- Command{Kind: kind, Args: args}
				}
			}
		}
	}
}

func mustJSON(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}
