package main

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/cmd/api/server"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/cmd/api/handlers"
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
	db := database.NewChatDB()
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handler := handlers.NewHandler(services)

	serv := new(server.ChatServer)

	go func() {
		if err := serv.Run("8080", handler.InitRoutes()); err != nil {
			log.Fatalf("error running chat server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	<-stop

	if err := serv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error shutting down server: %s", err)
	}
}
