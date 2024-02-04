package main

import (
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/cmd/api/handlers"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/chatutil"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"log"
)

// @title        Chat API
// @version      1.0
// @description  API Server for Chat Application

// @host localhost:8080
// @BasePath  /

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
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
	if err := server.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error running chat server: %s", err)
	}
}
