package handler

import (
	"fmt"
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

type UserListsResponse struct {
	Status string         `json:"status"`
	Result []*models.List `json:"result"`
}

type ItemCreateResponse struct {
	Status string `json:"status"`
	ItemID int64  `json:"item_id"`
}

type ListCreateResponse struct {
	Status string `json:"status"`
	ListID int64  `json:"list_id"`
}

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

func (h *Handler) GetUserId(c *gin.Context) (int64, error) {
	uid, ok := c.Get(idCtx)
	if !ok {
		return 0, fmt.Errorf("can't get userID from context")
	}

	var res int64
	if res, ok = uid.(int64); !ok {
		return 0, fmt.Errorf("can't convert interface userID to int64")
	}

	return res, nil
}

func (h *Handler) GetUserUUID(c *gin.Context) (string, error) {
	uUUID, ok := c.Get(CtxUUID)
	if !ok {
		return "", fmt.Errorf("can't get userUUID from context")
	}

	var res string
	if res, ok = uUUID.(string); !ok {
		return "", fmt.Errorf("can't convert interface userUUID to int64")
	}

	return res, nil
}
