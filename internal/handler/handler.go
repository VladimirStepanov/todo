package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	UserService  models.UserService
	MailService  models.MailService
	TokenService models.TokenService
	logger       *logrus.Logger
}

func (h *Handler) InitRoutes(mode string) http.Handler {
	gin.SetMode(mode)
	r := gin.New()

	auth := r.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.GET("/confirm/:link", h.confirm)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/refresh", h.refreshToken)
		auth.GET("/logout", h.authMiddleware, h.logout)
	}

	return r
}

func New(
	UserService models.UserService,
	MailService models.MailService,
	TokenService models.TokenService,
	logger *logrus.Logger) *Handler {

	return &Handler{
		UserService:  UserService,
		MailService:  MailService,
		TokenService: TokenService,
		logger:       logger}
}
