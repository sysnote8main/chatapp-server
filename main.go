package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

type Msg struct {
	Username string
	Message  string
}

var (
	clients      = make(map[*websocket.Conn]bool)
	msgBroadcast = make(chan Msg)
	upgrader     = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func broadcastMessages() {
	for {
		msg := <-msgBroadcast
		for c := range clients {
			err := c.WriteJSON(msg)
			if err != nil {
				slog.Error("Failed to write json", slog.Any("error", err))
				c.Close()
				delete(clients, c)
			}
		}
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade request", slog.Any("error", err))
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var data Msg
		err := conn.ReadJSON(&data)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Info("Connection Closed.")
			} else {
				slog.Error("Failed to read Json", slog.Any("error", err))
			}
			break
		}

		slog.Info("Message received!", slog.String("username", data.Username), slog.String("message", string(data.Message)))
		msgBroadcast <- data
	}
}

func main() {
	http.HandleFunc("/ws", handle)
	go broadcastMessages()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to listen http server", slog.Any("error", err))
	}
}
