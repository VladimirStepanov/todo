package handler

import (
	"net/http"
	"strconv"

	"github.com/VladimirStepanov/todo-app/internal/models"
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

	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)
	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrBadParam.Error(),
		})
		return
	}

	item, err := h.ItemService.GetItemByID(listID, itemID)

	if err != nil {
		switch err {
		case models.ErrNoItem:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, item)

}

func (h *Handler) updateItem(c *gin.Context) {
	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)
	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrBadParam.Error(),
		})
		return
	}

	req := &models.UpdateItemReq{}

	if !bindData(c, req) {
		return
	}

	err = h.ItemService.Update(listID, itemID, req)

	if err != nil {
		switch err {
		case models.ErrUpdateEmptyArgs, models.ErrTitleTooShort:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case models.ErrNoItem:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}

func (h *Handler) doneItem(c *gin.Context) {
}

func (h *Handler) deleteItem(c *gin.Context) {
	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)
	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrBadParam.Error(),
		})
		return
	}

	err = h.ItemService.Delete(listID, itemID)

	if err != nil {
		switch err {
		case models.ErrNoItem:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
