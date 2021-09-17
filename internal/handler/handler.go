package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Handler struct {
	UserService  models.UserService
	MailService  models.MailService
	TokenService models.TokenService
	ListService  models.ListService
	ItemService  models.ItemService
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (h *Handler) InitRoutes(mode string) http.Handler {
	gin.SetMode(mode)
	r := gin.New()

	r.Use(CORSMiddleware())
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
			lists.PATCH("/:list_id", h.onlyAdminAccessMiddleware, h.updateList)
			lists.PATCH("/:list_id/edit-role", h.onlyAdminAccessMiddleware, h.editRole)
			lists.DELETE("/:list_id", h.onlyAdminAccessMiddleware, h.deleteList)

			items := lists.Group("/:list_id/items", h.checkAccessToListMiddleware)
			{
				items.POST("", h.itemCreate)
				items.GET("", h.getItems)
				items.GET("/:item_id", h.getItemByID)
				items.PATCH("/:item_id", h.updateItem)
				items.PATCH("/:item_id/done", h.doneItem)
				items.DELETE("/:item_id", h.deleteItem)
			}
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

func New(
	UserService models.UserService,
	MailService models.MailService,
	TokenService models.TokenService,
	ListService models.ListService,
	ItemService models.ItemService,
	logger *logrus.Logger) *Handler {

	return &Handler{
		UserService:  UserService,
		MailService:  MailService,
		TokenService: TokenService,
		ListService:  ListService,
		ItemService:  ItemService,
		logger:       logger}
}
