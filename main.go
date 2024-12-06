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
		err = conn.WriteMessage(websocket.TextMessage, []byte("Server got a message!"))
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
