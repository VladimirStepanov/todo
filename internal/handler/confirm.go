package handler

import (
	"net/http"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

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
