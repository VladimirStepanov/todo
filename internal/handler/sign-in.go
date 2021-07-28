package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) signIn(c *gin.Context) {
	var req signupReq
	if !bindData(c, &req) {
		return
	}

	user, err := h.UserService.SignIn(req.Email, req.Password)

	if err != nil {
		switch err {
		case models.ErrBadUser:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case models.ErrUserNotActivated:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	td, err := h.TokenService.NewTokenPair(user.ID)
	if err != nil {
		switch err {
		case models.ErrMaxLoggedIn:
			c.JSON(http.StatusForbidden, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	})
}
