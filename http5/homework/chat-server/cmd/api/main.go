package main

import (
	"context"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/api/handlers"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/api"
	"github.com/ew0s/ewos-to-go-hw/http5/homework/chat-server/internal/pkg/httputils/server"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	db := database.NewChatDB()

	repos := repository.NewRepository(db)

	services := service.NewService(repos)

	authHandler := handlers.NewAuthHandler(services.Auth)
	publicChatHandler := handlers.NewPublicChatHandler(services.Chat)
	privateChatHandler := handlers.NewPrivateChatHandler(services.Chat)

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
