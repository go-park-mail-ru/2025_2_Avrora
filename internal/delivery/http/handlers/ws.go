package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // иначе будет CORS FAIL
}

type RoomMessage struct {
	Room      string `json:"room"`
	UserID    string `json:"userId"`
	UserEmail string `json:"userEmail"`
	UserName  string `json:"userName"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}

var Rooms = make(map[string]map[*websocket.Conn]bool)
var Broadcast = make(chan RoomMessage)

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	if room == "" {
		http.Error(w, "room required", http.StatusBadRequest)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer ws.Close()

	// создаём комнату, если нет
	if _, ok := Rooms[room]; !ok {
		Rooms[room] = make(map[*websocket.Conn]bool)
	}

	// добавляем юзера в комнату
	Rooms[room][ws] = true
	log.Printf("client joined room %s\n", room)

	for {
		var msg RoomMessage
		if err := ws.ReadJSON(&msg); err != nil {
			log.Println("ReadJSON error:", err)
			delete(Rooms[room], ws)
			log.Printf("client disconnected from %s\n", room)
			break
		}

		msg.Room = room // на всякий случай
		Broadcast <- msg
	}
}

// Рассылка
func HandleMessages() {
	for {
		msg := <-Broadcast

		// берем только клиентов из нужной комнаты
		clients := Rooms[msg.Room]

		for client := range clients {
			if err := client.WriteJSON(msg); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}
