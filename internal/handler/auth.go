package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

type signupReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=8,lte=32"`
}

func (h *Handler) signUp(c *gin.Context) {
	var req signupReq
	if ok := bindData(c, &req); !ok {
		return
	}

	user, err := h.UserService.Create(req.Email, req.Password)

	if err != nil {
		if err == models.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		} else {
			h.InternalError(c, err)
		}
		return
	}

	err = h.MailService.SendConfirmationsEmail(user)

	if err != nil {
		h.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (h *Handler) signIn(c *gin.Context) {
	c.String(http.StatusOK, "Sign in")
}

func (h *Handler) refreshToken(c *gin.Context) {
	c.String(http.StatusOK, "Refresh token")
}
