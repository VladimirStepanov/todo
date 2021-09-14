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
			c.JSON(http.StatusUnauthorized, gin.H{
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
			c.JSON(http.StatusUnprocessableEntity, gin.H{
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

func (h *Handler) confirm(c *gin.Context) {
	link := c.Param("link")

	err := h.UserService.ConfirmEmail(link)

	if err != nil {
		if err == models.ErrConfirmLinkNotExists {
			h.PageNotFound(c)
		} else {
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

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
			c.JSON(http.StatusBadRequest, gin.H{
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
