package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type BindDataError struct {
	Err         string            `json:"error"`
	InvalidArgs []invalidArgument `json:"invalidArgs"`
}

func (h *Handler) InternalError(c *gin.Context, err error) {
	h.logger.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": "Internal server error",
	})
}
