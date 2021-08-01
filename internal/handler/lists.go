package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type listCreateReq struct {
	Title       string `json:"title" binding:"required,gte=1,lte=255"`
	Description string `json:"description" binding:"required"`
}

func (h *Handler) listCreate(c *gin.Context) {
	var req listCreateReq
	if ok := bindData(c, &req); !ok {
		return
	}

	userID, err := h.GetUserId(c)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	listID, err := h.ListService.Create(req.Title, req.Description, userID)

	if err != nil {
		h.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"list_id": listID,
	})
}
