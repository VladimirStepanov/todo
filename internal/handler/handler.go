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
	ListService  models.ListService
	logger       *logrus.Logger
}

func (h *Handler) AccessLogger(c *gin.Context) {
	c.Next()

	h.logger.WithFields(logrus.Fields{
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
		"code":   c.Writer.Status(),
	}).Info("access")
}

func (h *Handler) InitRoutes(mode string) http.Handler {
	gin.SetMode(mode)
	r := gin.New()

	r.Use(h.AccessLogger)

	auth := r.Group("/auth")
	{
		auth.POST("/sign-in", h.signIn)
		auth.GET("/confirm/:link", h.confirm)
		auth.POST("/sign-up", h.signUp)
		auth.POST("/refresh", h.refreshToken)
		auth.GET("/logout", h.authMiddleware, h.logout)
	}

	api := r.Group("/api", h.authMiddleware)
	{
		lists := api.Group("/lists")
		{
			lists.POST("", h.listCreate)
			lists.GET("", h.getUserLists)
			lists.GET("/:list_id", h.getListByID)
			lists.PATCH("/:list_id", h.onlyAdminAccess, h.updateList)
			lists.PATCH("/:list_id/edit-role", h.onlyAdminAccess, h.editRole)
			lists.DELETE("/:list_id", h.onlyAdminAccess, h.deleteList)
		}
	}
	return r
}

func New(
	UserService models.UserService,
	MailService models.MailService,
	TokenService models.TokenService,
	ListService models.ListService,
	logger *logrus.Logger) *Handler {

	return &Handler{
		UserService:  UserService,
		MailService:  MailService,
		TokenService: TokenService,
		ListService:  ListService,
		logger:       logger}
}
