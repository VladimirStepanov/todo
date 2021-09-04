package handler

import (
	"net/http"
	"strconv"
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

func (h *Handler) checkAdminAccess(c *gin.Context) error {
	userID, err := h.GetUserId(c)
	if err != nil {
		return err
	}

	listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
	if err != nil {
		return models.ErrBadParam
	}

	return h.ListService.IsListAdmin(listID, userID)
}

func (h *Handler) onlyAdminAccessMiddleware(c *gin.Context) {
	err := h.checkAdminAccess(c)

	if err != nil {
		switch err {
		case models.ErrBadParam:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})

		case models.ErrNoList:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case models.ErrNoListAccess:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		c.Abort()
		return
	}

	c.Next()
}
