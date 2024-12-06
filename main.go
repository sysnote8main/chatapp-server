package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Failed to upgrade request", slog.Any("error", err))
	}
	defer conn.Close()

	for {
		msgType, msgByte, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message", slog.Any("error", err))
			break
		}

		slog.Info("Message received!", slog.String("message", string(msgByte)))
		err = conn.WriteMessage(msgType, []byte("Server got a message!"))
		if err != nil {
			slog.Error("Failed to respond message", slog.Any("error", err))
			break
		}
	}
}

func main() {
	http.HandleFunc("/ws", handle)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error("Failed to listen http server", slog.Any("error", err))
	}
}
