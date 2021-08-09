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
