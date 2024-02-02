package chatutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

type ChatServer struct {
	httpServer   *http.Server
	Users        map[string]string
	UserMutex    sync.Mutex
	ChatMessages []Message
	MessageMutex sync.Mutex
	UserMessages map[string][]Message
	UserMsgMutex sync.Mutex
}

func (cs *ChatServer) Run(port string, handler http.Handler) error {
	cs.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return cs.httpServer.ListenAndServe()
}

func (cs *ChatServer) Shutdown(ctx context.Context) error {
	return cs.httpServer.Shutdown(ctx)
}

func (cs *ChatServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cs.UserMutex.Lock()
	defer cs.UserMutex.Unlock()

	if _, exists := cs.Users[newUser.Username]; exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	cs.Users[newUser.Username] = newUser.Password
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s registered successfully", newUser.Username)
}

func (cs *ChatServer) SendMessage(w http.ResponseWriter, r *http.Request) {
	var newMessage Message
	err := json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cs.MessageMutex.Lock()
	defer cs.MessageMutex.Unlock()

	cs.ChatMessages = append(cs.ChatMessages, newMessage)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Message from %s sent successfully", newMessage.Username)
}
