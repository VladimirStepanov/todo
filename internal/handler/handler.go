package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	UserService models.UserService
	MailService models.MailService
	logger      *logrus.Logger
}

func (h *Handler) InitRoutes(mode string) http.Handler {
	gin.SetMode(mode)
	r := gin.New()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-up", h.signUp)
		auth.GET("/refresh", h.refreshToken)
	}

	return r
}

func New(UserService models.UserService, MailService models.MailService, logger *logrus.Logger) *Handler {
	return &Handler{UserService: UserService, MailService: MailService, logger: logger}
}
