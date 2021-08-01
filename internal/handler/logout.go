package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) logout(c *gin.Context) {
	userID, err := h.GetUserId(c)
	if err != nil {
		h.InternalError(c, err)
		return
	}
	userUUID, err := h.GetUserUUID(c)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	err = h.TokenService.Logout(userID, userUUID)
	if err != nil {
		h.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
