package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

type refreshReq struct {
	Token string `json:"refresh_token" binding:"required"`
}

func (h *Handler) refreshToken(c *gin.Context) {
	var req refreshReq
	if ok := bindData(c, &req); !ok {
		return
	}

	td, err := h.TokenService.Refresh(req.Token)

	if err != nil {
		switch err {
		case models.ErrBadToken:
			c.JSON(http.StatusForbidden, gin.H{
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
		return
	}

	c.JSON(http.StatusOK, &TokensResponse{
		AccessToken:  td.AccessToken,
		RefreshToken: td.RefreshToken,
	})
}
