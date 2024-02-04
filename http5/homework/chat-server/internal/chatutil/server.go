package chatutil

import (
	"context"
	"net/http"
	"time"
)

const (
	maxHeaderBytes = 1 << 20
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
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
		MaxHeaderBytes: maxHeaderBytes,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
	}

	return cs.httpServer.ListenAndServe()
}

func (cs *ChatServer) Shutdown(ctx context.Context) error {
	return cs.httpServer.Shutdown(ctx)
}
