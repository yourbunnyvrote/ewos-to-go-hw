package chatutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// User представляет пользователя чата
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Message представляет сообщение в чате
type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

// ChatServer представляет веб-сервер чата
type ChatServer struct {
	httpServer   *http.Server
	Users        map[string]string    // Хранение пользователей (username: password)
	UserMutex    sync.Mutex           // Мьютекс для безопасного доступа к Users
	ChatMessages []Message            // Хранение сообщений чата
	MessageMutex sync.Mutex           // Мьютекс для безопасного доступа к ChatMessages
	UserMessages map[string][]Message // Хранение личных сообщений (username: messages)
	UserMsgMutex sync.Mutex           // Мьютекс для безопасного доступа к UserMessages
}

func (cs *ChatServer) Run(port string, handler http.Handler) error {
	cs.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	return cs.httpServer.ListenAndServe()
}

func (cs *ChatServer) Shutdown(ctx context.Context) error {
	return cs.httpServer.Shutdown(ctx)
}

// RegisterUser представляет регистрацию нового пользователя
func (cs *ChatServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	cs.UserMutex.Lock()
	defer cs.UserMutex.Unlock()

	// Проверка, что пользователь с таким именем не существует
	if _, exists := cs.Users[newUser.Username]; exists {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Добавление нового пользователя
	cs.Users[newUser.Username] = newUser.Password
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s registered successfully", newUser.Username)
}

// SendMessage представляет отправку сообщения в общий чат
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
