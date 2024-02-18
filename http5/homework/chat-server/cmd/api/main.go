package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/ew0s/ewos-to-go-hw/docs"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/middleware"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/public_message"
	"github.com/ew0s/ewos-to-go-hw/internal/database"
	auth2 "github.com/ew0s/ewos-to-go-hw/internal/repository/auth"
	privateMessage2 "github.com/ew0s/ewos-to-go-hw/internal/repository/private_message"
	publicMessage2 "github.com/ew0s/ewos-to-go-hw/internal/repository/public_message"
	auth3 "github.com/ew0s/ewos-to-go-hw/internal/service/auth"
	privateMessage3 "github.com/ew0s/ewos-to-go-hw/internal/service/private_message"
	publicMessage3 "github.com/ew0s/ewos-to-go-hw/internal/service/public_message"
	"github.com/ew0s/ewos-to-go-hw/pkg/api"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/server"

	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Chat API
//	@version		1.0
//	@description	API Server for Chat Application

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.basic	BasicAuth
// @in							header
// @name						Authorization
func main() {
	chatDB := database.NewChatDB()

	authRepo := auth2.NewRepository(chatDB)
	privateMessageRepo := privateMessage2.NewRepository(chatDB)
	publicMessageRepo := publicMessage2.NewRepository(chatDB)

	authService := auth3.NewService(authRepo)
	privateMessageService := privateMessage3.NewService(privateMessageRepo)
	publicMessageService := publicMessage3.NewService(publicMessageRepo)

	authHandler := auth.NewHandler(authService)
	userIdentity := middleware.NewUserIdentity(authService)
	privateChatHandler := private_message.NewPrivateChatHandler(privateMessageService, userIdentity)
	publicChatHandler := public_message.NewPublicChatHandler(publicMessageService, userIdentity)

	routers := map[string]chi.Router{
		AuthEndpoint:           authHandler.Routes(),
		PublicMessageEndpoint:  publicChatHandler.Routes(),
		PrivateMessageEndpoint: privateChatHandler.Routes(),
	}

	r := api.MakeRoutes(AppVersion, routers)

	r.Get(SwaggerEndpoint, httpSwagger.Handler(
		httpSwagger.URL(DocJSONPath),
	))

	serv := new(server.ChatServer)

	go func() {
		if err := serv.Run(ServerPort, r); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error running chat server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	<-stop

	if err := serv.Shutdown(context.Background()); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error shutting down server: %s", err)
		}

		log.Println("server shut down gracefully")
	}
}
