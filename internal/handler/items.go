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

// CreateItem godoc
// @Summary Create item
// @Tags items
// @Accept  json
// @Produce  json
// @ID create-item
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param input body itemCreateReq true "item input"
// @Success 200 {object} ItemCreateResponse "success item creation"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/items [post]
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

// GetItems godoc
// @Summary Get items
// @Tags items
// @Accept  json
// @Produce  json
// @ID get-items
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Success 200 {object} UserItemsResponse "all user items"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/items [get]
func (h *Handler) getItems(c *gin.Context) {
	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)

	result, err := h.ItemService.GetItems(listID)

	if err != nil {
		h.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "result": result})

}

// GetItem godoc
// @Summary Get item
// @Tags items
// @Accept  json
// @Produce  json
// @ID get-item
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param item_id path int true "item_id"
// @Success 200 {object} models.Item "item"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found, item not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/items/{item_id} [get]
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

// UpdateItem godoc
// @Summary Update item
// @Tags items
// @Accept  json
// @Produce  json
// @ID update-item
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param item_id path int true "item_id"
// @Param input body models.UpdateItemReq true "input"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found, item not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/items/{item_id} [patch]
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

// DoneItem godoc
// @Summary Done item
// @Tags items
// @Accept  json
// @Produce  json
// @ID done-item
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param item_id path int true "item_id"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found, item not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/items/{item_id}/done [patch]
func (h *Handler) doneItem(c *gin.Context) {
	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)
	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrBadParam.Error(),
		})
		return
	}

	err = h.ItemService.Done(listID, itemID)

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

// DeleteItem godoc
// @Summary Delete item
// @Tags items
// @Accept  json
// @Produce  json
// @ID delete-item
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param item_id path int true "item_id"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found, item not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/items/{item_id} [delete]
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
