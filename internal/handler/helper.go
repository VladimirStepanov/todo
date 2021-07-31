package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BindDataError struct {
	Err         string            `json:"error"`
	InvalidArgs []invalidArgument `json:"invalidArgs"`
}

func (h *Handler) PageNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "Page not found",
	})
}

func (h *Handler) InternalError(c *gin.Context, err error) {
	h.logger.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": "Internal server error",
	})
}
