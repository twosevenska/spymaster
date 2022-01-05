package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"spymaster/src/spymaster"
	"spymaster/types"
)

var partialSearchFields = []string{"first_name", "last_name", "nickname", "email"}
var exactSearchFields = []string{"_id", "country"}

// ListUsers lists all users that meet the query parameters
func ListUsers(c *gin.Context) {
	perPage := c.MustGet("per_page").(int)
	pageNumber := c.MustGet("page_number").(int)

	exact := map[string]interface{}{}
	for _, field := range exactSearchFields {
		value, found := c.GetQuery(field)
		if found {
			exact[field] = value
		}
	}

	partial := map[string]string{}
	for _, field := range partialSearchFields {
		value, found := c.GetQuery(field)
		if found {
			partial[field] = value
		}
	}

	response, err := spymaster.ListUsers(c, exact, partial, perPage, pageNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	WritePaginationHeaders(c, response.TotalCount)
	c.JSON(http.StatusOK, response)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	var payload = &types.UserPost{}
	errs := c.ShouldBindJSON(payload)
	if errs != nil {
		e := fmt.Sprintf("Invalid payload received: %s", errs)
		c.JSON(http.StatusBadRequest, e)
		return
	}

	user, err := spymaster.CreateUser(c, payload)
	if err != nil {
		if err == spymaster.ErrDup {
			c.JSON(http.StatusConflict, gin.H{"message": "User/Email already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// UpdateUser creates a new user
func UpdateUser(c *gin.Context) {
	id, _ := c.GetQuery("id")

	var payload = &types.UserPatch{}
	errs := c.ShouldBindWith(payload, binding.JSON)
	fmt.Printf("\n%#v\n", payload)
	fmt.Printf("\n%#v\n", &types.UserPatch{})
	if (errs != nil || *payload == types.UserPatch{}) {
		e := fmt.Sprintf("Invalid payload received: %s", errs)
		c.JSON(http.StatusBadRequest, e)
		return
	}

	user, err := spymaster.UpdateUser(c, id, payload)
	if err != nil {
		if err == spymaster.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser creates a new user
func DeleteUser(c *gin.Context) {
	id, _ := c.GetQuery("id")

	err := spymaster.DeleteUser(c, id)
	if err != nil {
		if err == spymaster.ErrNotFound {
			c.Status(http.StatusNoContent)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
