package main

import (
	"log"

	"github.com/VladimirStepanov/todo-app/internal/config"
	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/repository/postgres"
	"github.com/VladimirStepanov/todo-app/internal/repository/redisrepo"
	"github.com/VladimirStepanov/todo-app/internal/server"
	"github.com/VladimirStepanov/todo-app/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg, err := config.New(".env")

	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.NewDB(
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser,
		cfg.PostgresPass, cfg.PostgresDB, "disable",
	)

	if err != nil {
		log.Println("Can't create new database", err)
		return
	}

	redisClient, err := redisrepo.NewRedisClient(cfg.RedisHost, cfg.RedisPort)
	if err != nil {
		log.Println("Can't create new redis client", err)
		return
	}

	listRepo := postgres.NewPostgresListRepository(db)
	tokenRepo := redisrepo.NewRedisRepository(redisClient)
	userRepo := postgres.NewPostgresUserRepository(db)
	userService := service.NewUserService(userRepo)
	mailService := service.NewMailService(cfg.Email, cfg.EmailPassword, cfg.Domain)
	listService := service.NewListService(listRepo)
	tokenService := service.NewTokenService(
		cfg.AccessKey, cfg.RefreshKey,
		cfg.MaxLoggedIn, tokenRepo,
	)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	handler := handler.New(userService, mailService, tokenService, listService, logger)

	srv := server.New(cfg.GetServerAddr(), handler.InitRoutes(cfg.Mode))
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
