package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// used to help extract validation errors
type invalidArgument struct {
	Field string `json:"field"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
	Param string `json:"param"`
}

// bindData is helper function, returns false if data is not bound
func bindData(c *gin.Context, req interface{}) bool {
	if c.ContentType() != "application/json" {
		msg := fmt.Sprintf("%s only accepts Content-Type application/json", c.FullPath())

		c.JSON(http.StatusBadRequest, gin.H{
			"error": msg,
		})
		return false
	}
	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBindJSON(req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			// could probably extract this, it is also in middleware_auth_user
			var invalidArgs []invalidArgument

			for _, err := range errs {
				invalidArgs = append(invalidArgs, invalidArgument{
					err.Field(),
					fmt.Sprintf("%v", err.Value()),
					err.Tag(),
					err.Param(),
				})
			}

			msg := "Invalid request parameters. See invalidArgs"

			c.JSON(http.StatusBadRequest, gin.H{
				"error":       msg,
				"invalidArgs": invalidArgs,
			})
			return false
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "Can't parse JSON request"})
		return false
	}

	return true
}
