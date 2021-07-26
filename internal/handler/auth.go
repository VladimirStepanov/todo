package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) refreshToken(c *gin.Context) {
	c.String(http.StatusOK, "Refresh token")
}
