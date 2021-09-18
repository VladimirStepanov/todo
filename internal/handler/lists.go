package handler

import (
	"net/http"
	"strconv"

	"github.com/VladimirStepanov/todo-app/internal/models"
	"github.com/gin-gonic/gin"
)

type listCreateReq struct {
	Title       string `json:"title" binding:"required,gte=1,lte=255"`
	Description string `json:"description" binding:"required"`
}

// CreateList godoc
// @Summary Create list
// @Tags lists
// @Accept  json
// @Produce  json
// @ID create-list
// @Security ApiKeyAuth
// @Param input body listCreateReq true "list input"
// @Success 200 {object} ListCreateResponse "success list creation"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists [post]
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

	c.JSON(http.StatusOK, &ListCreateResponse{"success", listID})
}

// GetList godoc
// @Summary Get list by id
// @Tags lists
// @Produce  json
// @ID get-list
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Success 200 {object} models.List "list"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "list not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id} [get]
func (h *Handler) getListByID(c *gin.Context) {
	userID, err := h.GetUserId(c)
	if err != nil {
		h.InternalError(c, err)
		return
	}

	listID, err := strconv.ParseInt(c.Param("list_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": models.ErrBadParam.Error(),
		})
		return
	}

	userList, err := h.ListService.GetListByID(listID, userID)
	if err != nil {
		switch err {
		case models.ErrNoList:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, userList)
}

type editRoleReq struct {
	UserID  int64 `json:"user_id" binding:"required"`
	IsAdmin *bool `json:"is_admin" binding:"required"`
}

// EditRole godoc
// @Summary Edit user role for list
// @Tags lists
// @Accept  json
// @Produce  json
// @ID edit-role
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param input body editRoleReq true "edit-role input"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 403 {object} ErrorResponse "current user is not admin"
// @Failure 404 {object} ErrorResponse "list not found, user not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id}/edit-role [patch]
func (h *Handler) editRole(c *gin.Context) {
	var req editRoleReq
	if ok := bindData(c, &req); !ok {
		return
	}

	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)

	err := h.ListService.EditRole(listID, req.UserID, *(req.IsAdmin))

	if err != nil {
		switch err {
		case models.ErrUserNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// DeleteList godoc
// @Summary Delete list by id
// @Tags lists
// @Produce  json
// @ID delete-list
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 403 {object} ErrorResponse "current user is not admin"
// @Failure 404 {object} ErrorResponse "list not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)

	err := h.ListService.Delete(listID)

	if err != nil {
		switch err {
		case models.ErrNoList:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// UpdateList godoc
// @Summary Update list by id
// @Tags lists
// @Produce  json
// @ID delete-list
// @Security ApiKeyAuth
// @Param list_id path int true "list_id"
// @Param input body models.UpdateListReq true "input"
// @Success 200 {string} status	"success"
// @Failure 400 {object} ErrorResponse	"bad input, auth header errors"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 403 {object} ErrorResponse "current user is not admin"
// @Failure 404 {object} ErrorResponse "list not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists/{list_id} [patch]
func (h *Handler) updateList(c *gin.Context) {

	var req models.UpdateListReq
	if ok := bindData(c, &req); !ok {
		return
	}

	listID, _ := strconv.ParseInt(c.Param("list_id"), 10, 64)

	err := h.ListService.Update(listID, &req)

	if err != nil {
		switch err {
		case models.ErrUpdateEmptyArgs, models.ErrTitleTooShort:
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		case models.ErrNoList:
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
		default:
			h.InternalError(c, err)
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetUserLists godoc
// @Summary Get all user lists
// @Tags lists
// @Produce  json
// @ID get-lists
// @Security ApiKeyAuth
// @Success 200 {object} UserListsResponse "lists"
// @Failure 401 {object} ErrorResponse "user is not authorized"
// @Failure 404 {object} ErrorResponse "user not found"
// @Failure 500 {object} ErrorResponse "internal error"
// @Router /api/lists [get]
func (h *Handler) getUserLists(c *gin.Context) {
	result, err := h.ListService.GetUserLists(c.GetInt64(idCtx))

	if err != nil {
		h.InternalError(c, err)
		return
	}

	c.JSON(http.StatusOK, UserListsResponse{"success", result})
}
