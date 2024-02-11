package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/handlers"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/api"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/httputils/server"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/docs"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/database"

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
	chatDB := database.NewChatDB()

	chatRepo := repository.NewRepository(chatDB)

	chatService := service.NewService(chatRepo)

	userIdentity := handlers.NewUserIdentity(chatService.Auth)

	authHandler := handlers.NewAuthHandler(chatService.Auth)
	publicChatHandler := handlers.NewPublicChatHandler(chatService.Chat, userIdentity)
	privateChatHandler := handlers.NewPrivateChatHandler(chatService.Chat, userIdentity)

	routers := map[string]chi.Router{
		"/auth":             authHandler.Routes(),
		"/messages/public":  publicChatHandler.Routes(),
		"/messages/private": privateChatHandler.Routes(),
	}

	r := api.MakeRoutes("/v1", routers)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	serv := new(server.ChatServer)

	go func() {
		if err := serv.Run("8080", r); err != nil {
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
