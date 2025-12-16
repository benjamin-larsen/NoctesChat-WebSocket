package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/benjamin-larsen/NoctesChat-WebSocket/models"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write([]byte("{\"error\":\"Invalid WebSocket request\"}"))
	},
}

type Socket struct {
	conn      *websocket.Conn
	userId    uint64
	userToken models.UserToken
	hasAuth   bool
}

func (s *Socket) Close(closeCode int, text string) {
	s.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(closeCode, text),
		time.Now().Add(10*time.Second),
	)

	s.conn.Close()
}

func (s *Socket) Run() {
	for {
		_, rawMsg, err := s.conn.ReadMessage()

		if err != nil {
			return
		}

		var baseMsg models.BaseInbound
		err = json.Unmarshal(rawMsg, &baseMsg)

		if err != nil {
			s.Close(1008, "Invalid JSON")
			return
		}

		switch baseMsg.MsgType {
		case "login":
			{
				var msg models.LoginInbound
				err = json.Unmarshal(baseMsg.Data, &msg)

				if err != nil {
					s.Close(1008, "Invalid JSON")
					return
				}

				s.ProcessLogin(msg)
			}
		default:
			{
				s.Close(1008, "Unknown message type: " + baseMsg.MsgType)
				return
			}
		}
	}
}

func (s *Socket) Cleanup() {

}

func handleUpgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	socket := Socket{
		conn:    conn,
		hasAuth: false,
	}

	defer socket.Cleanup()
	defer conn.Close()

	socket.Run()
}

func SetupWS() {
	http.HandleFunc("/ws", handleUpgrade)

	log.Fatal(http.ListenAndServe(":3030", nil))
}
