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
	configApp "github.com/ew0s/ewos-to-go-hw/cmd/config"
	_ "github.com/ew0s/ewos-to-go-hw/docs"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/auth/request"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/middleware"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/private_message"
	"github.com/ew0s/ewos-to-go-hw/internal/api/handlers/public_message"
	"github.com/ew0s/ewos-to-go-hw/internal/database"
	"github.com/ew0s/ewos-to-go-hw/internal/domain/entities"
	reposInMemoryPrivateMessage "github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory/private_message"
	reposInMemoryPublicMessage "github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory/public_message"
	reposInMemoryAuth "github.com/ew0s/ewos-to-go-hw/internal/repository/inmemory/user"
	reposPostgresPrivateMessage "github.com/ew0s/ewos-to-go-hw/internal/repository/postgres/private_message"
	reposPostgresPublicMessage "github.com/ew0s/ewos-to-go-hw/internal/repository/postgres/public_message"
	reposPostgresAuth "github.com/ew0s/ewos-to-go-hw/internal/repository/postgres/user"
	serviceAuth "github.com/ew0s/ewos-to-go-hw/internal/service/auth"
	servicePrivateMessage "github.com/ew0s/ewos-to-go-hw/internal/service/private_message"
	servicePublicMessage "github.com/ew0s/ewos-to-go-hw/internal/service/public_message"
	"github.com/ew0s/ewos-to-go-hw/pkg/api"
	"github.com/ew0s/ewos-to-go-hw/pkg/httputils/server"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v4/stdlib"
	httpSwagger "github.com/swaggo/http-swagger"
)

type DataBase interface {
	Insert(query string, _ interface{}) error
	Get(query string) (interface{}, error)
}

type AuthRepository interface {
	CreateUser(credentials entities.AuthCredentials) error
	GetUser(username string) (entities.AuthCredentials, error)
}

type PrivateMessageRepository interface {
	SendPrivateMessage(receiver string, msg entities.Message) error
	GetPrivateChat(chat entities.ChatMetadata) ([]entities.Message, error)
	GetUserList(username string) ([]string, error)
}

type PublicMessageRepository interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicChat() ([]entities.Message, error)
}

type AuthService interface {
	CreateUser(req request.RegistrationRequest) (interface{}, error)
	GetUser(username string) (entities.AuthCredentials, error)
	Identify(interface{}) (string, error)
}

type PrivateMessageService interface {
	SendPrivateMessage(receiver string, msg entities.Message) error
	GetPrivateMessages(chat entities.ChatMetadata, params entities.PaginateParam) ([]entities.Message, error)
	GetUsersWithMessage(receiver string) ([]string, error)
}

type PublicMessageService interface {
	SendPublicMessage(msg entities.Message) error
	GetPublicMessages(params entities.PaginateParam) ([]entities.Message, error)
}

type UserIdentity interface {
	Identify(next http.Handler) http.Handler
}

//	@title			Chat API
//	@version		1.0
//	@description	API Server for Chat Application

//	@host		localhost:8080
//	@BasePath	/

// @securityDefinitions.basic	BasicAuth
// @in							header
// @name						Authorization
func main() {
	var (
		config configApp.Config

		db DataBase

		authRepo           AuthRepository
		privateMessageRepo PrivateMessageRepository
		publicMessageRepo  PublicMessageRepository

		authService           AuthService
		privateMessageService PrivateMessageService
		publicMessageService  PublicMessageService

		userIdentity UserIdentity

		authHandler           *auth.Handler
		privateMessageHandler *private_message.Handler
		publicMessageHandler  *public_message.Handler
	)

	err := config.ParseData(consts.ConfigPath)
	if err != nil {
		panic(err)
	}

	validate := validator.New()

	if config.StorageType == "postgres" {
		db = database.NewPostgresDB()

		authRepo = reposPostgresAuth.NewRepository(db)
		privateMessageRepo = reposPostgresPrivateMessage.NewRepository(db)
		publicMessageRepo = reposPostgresPublicMessage.NewRepository(db)

		privateMessageService = servicePrivateMessage.NewPostgresService(privateMessageRepo)
	} else {
		db = database.NewChatDB()

		authRepo = reposInMemoryAuth.NewRepository(db)
		privateMessageRepo = reposInMemoryPrivateMessage.NewRepository(db)
		publicMessageRepo = reposInMemoryPublicMessage.NewRepository(db)

		privateMessageService = servicePrivateMessage.NewService(privateMessageRepo)
	}

	publicMessageService = servicePublicMessage.NewService(publicMessageRepo)

	if config.AuthType == "jwt" {
		authService = serviceAuth.NewJWTService(authRepo)
		userIdentity = middleware.NewJWTIdentity(authService)
	} else {
		authService = serviceAuth.NewBasicAuthService(authRepo)
		userIdentity = middleware.NewUserIdentity(authService, validate)
	}

	authHandler = auth.NewHandler(authService, validate)
	privateMessageHandler = private_message.NewHandler(privateMessageService, userIdentity, validate)
	publicMessageHandler = public_message.NewHandler(publicMessageService, userIdentity, validate)

	routers := map[string]chi.Router{
		consts.AuthEndpoint:           authHandler.Routes(),
		consts.PublicMessageEndpoint:  publicMessageHandler.Routes(),
		consts.PrivateMessageEndpoint: privateMessageHandler.Routes(),
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
