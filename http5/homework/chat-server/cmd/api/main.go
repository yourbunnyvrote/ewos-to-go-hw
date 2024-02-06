package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/cmd/api/handlers"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
)

//	@title			Chat API
//	@version		1.0
//	@description	API Server for Chat Application

//	@host		localhost:8080
//	@BasePath	/v1

// @securityDefinitions.basic	BasicAuth
// @in							header
// @name						Authorization
func main() {
	db := chatutil.ChatDB{
		Users:        map[string]chatutil.User{},
		PublicChat:   []chatutil.Message{},
		PrivateChats: map[chatutil.Chat][]chatutil.Message{},
	}
	repos := repository.NewRepository(&db)
	services := service.NewService(repos)
	handler := handlers.NewHandler(services)

	server := new(chatutil.ChatServer)

	go func() {
		if err := server.Run("8080", handler.InitRoutes()); err != nil {
			log.Fatalf("error running chat server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	<-stop

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting down server: %s", err)
	}
}
