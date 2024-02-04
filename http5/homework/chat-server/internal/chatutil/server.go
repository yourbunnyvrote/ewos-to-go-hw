package chatutil

import (
	"context"
	"net/http"
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

type TextMessage struct {
	Content string `json:"content"`
}

type ChatServer struct {
	httpServer *http.Server
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
