package ws

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"sync"

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
	conn         *websocket.Conn
	userId       uint64
	userToken    models.UserToken
	hasAuth      bool
	isClosing    bool
	closingMutex sync.RWMutex
}

// WARNING: unlocking is expected to be handled by the calling function
// Returns wether the socket is closing, and wether the lock should be released
func (s *Socket) IsClosing() (isClosing bool, shouldUnlock bool) {
	lockSuccess := s.closingMutex.TryRLock()

	// if lock failed that means the Writer from Cleanup is running, which means that the Socket is indeed closing.
	if lockSuccess == false {
		return true, false
	}

	return s.isClosing, true
}

// This is a helper method for IsClosing() simply insert the second boolean from IsClosing() here in a defer
func (s *Socket) ReleaseCleanupLock(should bool) {
	if should {
		s.closingMutex.RUnlock()
	}
}

func (s *Socket) Close(closeCode int, text string) {
	s.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(closeCode, text),
		time.Now().Add(10*time.Second),
	)

	s.conn.Close()
}

func (s *Socket) RaiseException(err error) {
	s.Close(1011, "Internal Server Error")

	log.Printf("An error occured in WebSocket: %s\n", err)
}

var ErrInvalidJson = errors.New("Invalid JSON")
var ErrInvalidMessageType = errors.New("Unknown message type")

func (s *Socket) RunAuth() error {
	for s.hasAuth != true {
		_, rawMsg, err := s.conn.ReadMessage()

		if err != nil {
			return err
		}

		var baseMsg models.BaseInbound
		err = json.Unmarshal(rawMsg, &baseMsg)

		if err != nil {
			s.Close(1008, "Invalid JSON")
			return ErrInvalidJson
		}

		if baseMsg.MsgType != "login" {
			s.Close(1008, "Unknown message type: "+baseMsg.MsgType)
			return ErrInvalidMessageType
		}

		var msg models.LoginInbound
		err = json.Unmarshal(baseMsg.Data, &msg)

		if err != nil {
			s.Close(1008, "Invalid JSON")
			return err
		}

		err = s.ProcessLogin(msg)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Socket) Run() {
	err := s.RunAuth();

	if err != nil || s.hasAuth != true {
		return
	}

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
		default:
			{
				s.Close(1008, "Unknown message type: "+baseMsg.MsgType)
				return
			}
		}
	}
}

func (s *Socket) Cleanup() {
	s.closingMutex.Lock()
	defer s.closingMutex.Unlock()

	// Already doing cleanup
	if s.isClosing {
		return
	}

	s.isClosing = true
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
