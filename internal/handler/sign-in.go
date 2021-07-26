package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) signIn(c *gin.Context) {
	c.String(http.StatusOK, "Sign in")
}
