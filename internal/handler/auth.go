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

// SignIn godoc
// @Tags auth
// @Summary Sign in
// @Accept  json
// @Produce  json
// @ID login
// @Param input body signupReq true "credentials"
// @Success 200 {object} TokensResponse	"tokens"
// @Failure 400 {object} ErrorResponse	"bad input"
// @Failure 401 {object} ErrorResponse "user not activated"
// @Failure 404 {object} ErrorResponse "user not found"
// @Failure 422 {object} ErrorResponse "max logged in users in one account"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /auth/sign-in [post]
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

// SignUp godoc
// @Summary Sign up
// @Tags auth
// @Accept  json
// @Produce  json
// @ID register
// @Param input body signupReq true "register"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input"
// @Failure 409 {object} ErrorResponse "user already exists"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /auth/sign-up [post]
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

// Logout godoc
// @Summary Log out
// @Tags auth
// @Accept  json
// @Produce  json
// @ID logout
// @Security ApiKeyAuth
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse "auth header errors"
// @Failure 401 {object} ErrorResponse "user unauthorized"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /auth/logout [get]
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

// Confirm godoc
// @Summary Confirm email
// @Tags auth
// @Accept  json
// @Produce  json
// @ID confirm
// @Param link path string true "link confirmation"
// @Success 200 {string} status	"success"
// @Failure 404 {object} ErrorResponse	"page not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /auth/confirm/{link} [get]
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

// RefreshToken godoc
// @Summary Refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @ID refresh-token
// @Param input body refreshReq true "refresh token"
// @Success 200 {object} TokensResponse	"tokens"
// @Failure 400 {object} ErrorResponse	"bad token"
// @Failure 401 {object} ErrorResponse "user unauthorized"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /auth/refresh [post]
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
