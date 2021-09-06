package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type itemCreateReq struct {
	Title       string `json:"title" binding:"required,gte=1,lte=255"`
	Description string `json:"description" binding:"required"`
}

func (h *Handler) itemCreate(c *gin.Context) {

	var req itemCreateReq
	if !bindData(c, &req) {
		return
	}

	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)
	itemID, err := h.ItemService.Create(req.Title, req.Description, listID)

	if err != nil {
		h.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"item_id": itemID,
	})
}

func (h *Handler) getItems(c *gin.Context) {

}

func (h *Handler) getItemByID(c *gin.Context) {

}

func (h *Handler) updateItem(c *gin.Context) {

}

func (h *Handler) doneItem(c *gin.Context) {

}

func (h *Handler) deleteItem(c *gin.Context) {

}
