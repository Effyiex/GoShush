package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type ChatMessageContent struct {
	Body []byte `json:"body"`
	Type string `json:"type"`
}

type ChatMessage struct {
	Content  []ChatMessageContent `json:"content"`
	Sender   uuid.UUID            `json:"sender"`
	AtTime   string               `json:"time"`
	Received []uuid.UUID          `json:"received"`
}

type ChatRoom struct {
	ID       uuid.UUID     `json:"id"`
	Title    string        `json:"title"`
	Members  []uuid.UUID   `json:"members"`
	Messages []ChatMessage `json:"messages"`
}

const CHATS_STORAGE = "chats.json"

type ChatManager struct {
	Chats     []ChatRoom    `json:"chats"`
	Broadcast []ChatMessage `json:"broadcast"`
	Server    *Server       `json:"-"`
}

func BootChatManager(dataFolder string) ChatManager {
	cm := ChatManager{
		Chats:     make([]ChatRoom, 0),
		Broadcast: make([]ChatMessage, 0),
	}
	cm.Load(dataFolder)
	return cm
}

func (cm *ChatManager) RegisterOn(srv Server) {
	cm.Server = &srv
	srv.HandleFunc("/api/chats", cm.QueryChats()).Methods("GET")
	srv.HandleFunc("/api/broadcast/query", cm.QueryBroadcast()).Methods("GET")
	srv.HandleFunc("/api/broadcast/send", cm.SendBroadcast()).Methods("POST")
}

type SendBroadcastInPacket struct {
	Content []ChatMessageContent `json:"content"`
	Sender  uuid.UUID            `json:"sender"`
	Session uuid.UUID            `json:"session"`
}

func (cm *ChatManager) SendBroadcast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var inPacket SendBroadcastInPacket
		err := json.NewDecoder(r.Body).Decode(&inPacket)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if cm.Server.Users.IsValidSession(inPacket.Sender, inPacket.Session) {
			cm.Broadcast = append(cm.Broadcast, ChatMessage{
				Content:  inPacket.Content,
				Sender:   inPacket.Sender,
				AtTime:   time.Now().Format("2006-01-02 15:04:05"),
				Received: make([]uuid.UUID, 0),
			})
			w.Write([]byte("success"))
		} else {
			w.Write([]byte("invalid session"))
		}

	}
}

type QueryBroadcastOutPacket struct {
	Messages []ChatMessage `json:"messages"`
}

func (cm *ChatManager) QueryBroadcast() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params = r.URL.Query()
		var user, err_user = uuid.Parse(params.Get("user"))
		var session, err_session = uuid.Parse(params.Get("session"))
		if err_user != nil || err_session != nil {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}

		if !cm.Server.Users.IsValidSession(user, session) {
			w.Write([]byte("invalid session"))
			return
		}

		var messages []ChatMessage
		var overwrite []ChatMessage
		for _, message := range cm.Broadcast {
			var received = false
			for _, receiver := range message.Received {
				if receiver == user {
					received = true
					break
				}
			}
			if len(message.Received) != len(cm.Server.Users.Users) {
				overwrite = append(overwrite, message)
			}
			message.Received = append(message.Received, user)
			if !received {
				messages = append(messages, message)
			}
		}
		cm.Broadcast = overwrite

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&messages)

	}
}

func (cm *ChatManager) QueryChats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var params = r.URL.Query()
		var user, err_user = uuid.Parse(params.Get("user"))
		var session, err_session = uuid.Parse(params.Get("session"))
		if err_user != nil || err_session != nil {
			http.Error(w, "invalid params", http.StatusBadRequest)
			return
		}

		if !cm.Server.Users.IsValidSession(user, session) {
			w.Write([]byte("invalid session"))
			return
		}

		var chats []ChatRoom
		for _, chat := range cm.Chats {
			for _, member := range chat.Members {
				if member == user {
					chats = append(chats, chat)
					break
				}
			}
		}

		json.NewEncoder(w).Encode(&chats)

	}
}

func (cm *ChatManager) Load(dataFolder string) {
	data, err := os.ReadFile(dataFolder + "/" + CHATS_STORAGE)
	if err != nil {
		cm.Save(dataFolder)
		return
	}
	err = json.Unmarshal(data, cm)
	if err != nil {
		panic(err)
	}
}

func (cm *ChatManager) Save(dataFolder string) {
	data, _ := json.MarshalIndent(cm, "", "  ")
	os.WriteFile(dataFolder+"/"+CHATS_STORAGE, data, 0644)
}
