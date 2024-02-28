package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ew0s/ewos-to-go-hw/cmd/api/consts"
	_ "github.com/ew0s/ewos-to-go-hw/docs"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/middleware"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/public_message"
	"github.com/ew0s/ewos-to-go-hw/internal/database"
	reposAuth "github.com/ew0s/ewos-to-go-hw/internal/repository/auth"
	reposPrivateMessage "github.com/ew0s/ewos-to-go-hw/internal/repository/private_message"
	reposPublicMessage "github.com/ew0s/ewos-to-go-hw/internal/repository/public_message"
	serviceAuth "github.com/ew0s/ewos-to-go-hw/internal/service/auth"
	servicePrivateMessage "github.com/ew0s/ewos-to-go-hw/internal/service/private_message"
	servicePublicMessage "github.com/ew0s/ewos-to-go-hw/internal/service/public_message"
	"github.com/ew0s/ewos-to-go-hw/pkg/api"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/server"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
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
	validate := validator.New()

	chatDB := database.NewChatDB()

	authRepo := reposAuth.NewRepository(chatDB)
	privateMessageRepo := reposPrivateMessage.NewRepository(chatDB)
	publicMessageRepo := reposPublicMessage.NewRepository(chatDB)

	authService := serviceAuth.NewService(authRepo)
	privateMessageService := servicePrivateMessage.NewService(privateMessageRepo)
	publicMessageService := servicePublicMessage.NewService(publicMessageRepo)

	authHandler := auth.NewAuthHandler(authService, validate)
	userIdentity := middleware.NewUserIdentity(authService)
	privateChatHandler := private_message.NewPrivateChatHandler(privateMessageService, userIdentity, validate)
	publicChatHandler := public_message.NewPublicChatHandler(publicMessageService, userIdentity, validate)

	routers := map[string]chi.Router{
		consts.AuthEndpoint:           authHandler.Routes(),
		consts.PublicMessageEndpoint:  publicChatHandler.Routes(),
		consts.PrivateMessageEndpoint: privateChatHandler.Routes(),
	}

	r := api.MakeRoutes(consts.AppVersion, routers)

	r.Get(consts.SwaggerEndpoint, httpSwagger.Handler(
		httpSwagger.URL(consts.DocJSONPath),
	))

	serv := new(server.ChatServer)

	go func() {
		if err := serv.Run(consts.ServerPort, r); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
