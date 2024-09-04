package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]string) // client connection -> username
var messages = make(chan Message)
var mu sync.Mutex

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
	To       string `json:"to,omitempty"` // recipient's username
}

func (h *Handler) handleConnection(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error while upgrading connection:", err)
		return
	}
	defer conn.Close()

	// Read username first
	_, username, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error while reading username:", err)
		return
	}
	usernameStr := string(username)

	mu.Lock()
	clients[conn] = usernameStr
	mu.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error while reading message:", err)
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			fmt.Println("Error while unmarshaling message:", err)
			continue
		}

		msg.Username = usernameStr
		messages <- msg
	}
}

func HandleMessages() {
	for msg := range messages {
		mu.Lock()
		for client, username := range clients {
			if msg.To == "" || msg.To == username {
				fmt.Println(msg)
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Println("Error while writing message:", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
		mu.Unlock()
	}
}
