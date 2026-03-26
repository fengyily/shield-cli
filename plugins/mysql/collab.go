package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ── WebSocket upgrader ──

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// ── Collab message types ──

type CollabMsg struct {
	Type   string          `json:"type"`
	UserID string          `json:"user_id,omitempty"`
	Name   string          `json:"name,omitempty"`
	Color  string          `json:"color,omitempty"`
	DB     string          `json:"db,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

type CursorData struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type DragData struct {
	Table string  `json:"table"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
	Done  bool    `json:"done,omitempty"`
}

type PresenceUser struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}

// ── Client ──

type Client struct {
	hub    *CollabHub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	name   string
	color  string
	db     string // which database's ER diagram the user is viewing
}

// ── Hub ──

type CollabHub struct {
	mu         sync.RWMutex
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func newCollabHub() *CollabHub {
	return &CollabHub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte, 256),
	}
}

func (h *CollabHub) run() {
	for {
		select {
		case c := <-h.register:
			h.mu.Lock()
			h.clients[c] = true
			h.mu.Unlock()
			h.broadcastPresence(c.db)

		case c := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
			h.mu.Unlock()
			h.broadcastPresence(c.db)

		case msg := <-h.broadcast:
			h.mu.RLock()
			for c := range h.clients {
				select {
				case c.send <- msg:
				default:
					// drop slow client
				}
			}
			h.mu.RUnlock()
		}
	}
}

// broadcastPresence sends the user list to all clients viewing the given database.
func (h *CollabHub) broadcastPresence(db string) {
	h.mu.RLock()
	var users []PresenceUser
	var targets []*Client
	for c := range h.clients {
		if c.db == db {
			users = append(users, PresenceUser{UserID: c.userID, Name: c.name, Color: c.color})
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()

	data, _ := json.Marshal(users)
	msg := CollabMsg{Type: "presence", DB: db, Data: data}
	raw, _ := json.Marshal(msg)

	for _, c := range targets {
		select {
		case c.send <- raw:
		default:
		}
	}
}

// broadcastToDB sends a message to all clients viewing the database, except the sender.
func (h *CollabHub) broadcastToDB(db string, sender *Client, msg []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.db == db && c != sender {
			select {
			case c.send <- msg:
			default:
			}
		}
	}
}

// broadcastDBChanged notifies all clients viewing a database that the schema structure changed.
func (h *CollabHub) broadcastDBChanged(db string, sender *Client) {
	msg := CollabMsg{Type: "schema_changed", DB: db}
	raw, _ := json.Marshal(msg)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		if c.db == db && c != sender {
			select {
			case c.send <- raw:
			default:
			}
		}
	}
}

// ── Random name / color generation ──

var userNames = []string{
	"Panda", "Fox", "Eagle", "Wolf", "Bear", "Tiger", "Hawk", "Owl",
	"Lynx", "Otter", "Raven", "Crane", "Dolphin", "Falcon", "Koala",
	"Heron", "Parrot", "Gecko", "Bison", "Elk",
}

var userColors = []string{
	"#ef4444", "#f97316", "#eab308", "#22c55e", "#06b6d4",
	"#3b82f6", "#8b5cf6", "#ec4899", "#14b8a6", "#f43f5e",
}

func randomName() string {
	return userNames[rand.Intn(len(userNames))]
}

func randomColor() string {
	return userColors[rand.Intn(len(userColors))]
}

func randomID() string {
	return fmt.Sprintf("%d%04d", time.Now().UnixMilli()%100000, rand.Intn(10000))
}

// ── WebSocket handler ──

func collabHandler(hub *CollabHub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		db := r.URL.Query().Get("db")

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 64),
			userID: randomID(),
			name:   randomName(),
			color:  randomColor(),
			db:     db,
		}

		hub.register <- client

		// Send the client its own identity
		welcome := CollabMsg{
			Type:   "welcome",
			UserID: client.userID,
			Name:   client.name,
			Color:  client.color,
			DB:     db,
		}
		if raw, err := json.Marshal(welcome); err == nil {
			client.send <- raw
		}

		go client.writePump()
		go client.readPump()
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(4096)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg CollabMsg
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}

		// Stamp sender info
		msg.UserID = c.userID
		msg.Name = c.name
		msg.Color = c.color
		msg.DB = c.db

		switch msg.Type {
		case "cursor", "drag":
			// Forward to other clients in the same database
			stamped, _ := json.Marshal(msg)
			c.hub.broadcastToDB(c.db, c, stamped)

		case "schema_changed":
			c.hub.broadcastDBChanged(c.db, c)

		case "switch_db":
			// Client switched to viewing a different database
			oldDB := c.db
			var data struct {
				DB string `json:"db"`
			}
			if json.Unmarshal(msg.Data, &data) == nil && data.DB != "" {
				c.hub.mu.Lock()
				c.db = data.DB
				c.hub.mu.Unlock()
				c.hub.broadcastPresence(oldDB)
				c.hub.broadcastPresence(c.db)
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
