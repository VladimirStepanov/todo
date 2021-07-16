package main

import (
	"log"

	"github.com/VladimirStepanov/todo-app/internal/config"
	"github.com/VladimirStepanov/todo-app/internal/handler"
	"github.com/VladimirStepanov/todo-app/internal/repository/postgres"
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

	userRepo := postgres.NewPostgresRepository(db)
	userService := service.NewUserService(userRepo)
	mailService := service.NewMailService(cfg.Email, cfg.EmailPassword, cfg.Domain)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	handler := handler.New(userService, mailService, logger)

	srv := server.New(cfg.GetServerAddr(), handler.InitRoutes(cfg.Mode))
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
