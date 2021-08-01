package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) logout(c *gin.Context) {
	userID := c.GetInt64(idCtx)
	userUUID := c.GetString(CtxUUID)
	err := h.TokenService.Logout(userID, userUUID)
	if err != nil {
		h.InternalError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
