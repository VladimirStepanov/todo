package handler

import (
	"net/http"
	"strings"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

var (
	idCtx   = "CtxUserID"
	CtxUUID = "CtxUUID"
)

func (h *Handler) authMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrNoAuthHeader.Error(),
		})
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrInvalidAuthHeader.Error(),
		})
		c.Abort()
		return
	}

	userID, userUUID, err := h.TokenService.Verify(headerParts[1])

	if err != nil {
		switch err {
		case models.ErrBadToken:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case models.ErrUserUnauthorized, models.ErrTokenExpired:
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		c.Abort()
		return
	}

	c.Set(idCtx, userID)
	c.Set(CtxUUID, userUUID)
	c.Next()
}
