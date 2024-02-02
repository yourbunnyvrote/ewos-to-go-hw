package main

import (
	chatutil "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/handlers"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/repository"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/service"
	"log"
)

func main() {
	db := chatutil.ChatDB{
		Users:      map[string]chatutil.User{},
		PublicChat: []chatutil.Message{},
	}
	repos := repository.NewRepository(&db)
	services := service.NewService(repos)
	handler := handlers.NewHandler(services)

	server := new(chatutil.ChatServer)
	if err := server.Run("8080", handler.InitRoutes()); err != nil {
		log.Fatalf("error running chat server: %s", err)
	}
}
