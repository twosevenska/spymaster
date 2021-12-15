package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"internal/mongo"
	"types"
)

var partialSearchFields = []string{"first_name", "last_name", "nickname", "email"}
var exactSearchFields = []string{"_id", "country"}

// ListUsers lists all users that meet the query parameters
func ListUsers(c *gin.Context) {
	mc := c.MustGet("mongo").(mongo.Client)

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

	users, totalCount, err := mc.ListUsers(exact, partial, perPage, pageNumber)
	if err != nil {
		log.Printf("ListUsers: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Server error"})
		return
	}

	response := types.UsersResult{
		Page:       pageNumber,
		PerPage:    perPage,
		TotalCount: totalCount,
		Users:      users,
	}

	WritePaginationHeaders(c, totalCount)
	c.JSON(http.StatusOK, response)
}

// CreateUser creates a new user
func CreateUser(c *gin.Context) {
	mc := c.MustGet("mongo").(mongo.Client)

	var payload = &types.UserPost{}
	errs := c.ShouldBindJSON(payload)
	if errs != nil {
		e := fmt.Sprintf("Invalid payload received: %s", errs)
		c.JSON(http.StatusBadRequest, e)
		return
	}

	user, err := mc.CreateUser(*payload)
	if err != nil {
		if mc.IsDup(err) {
			c.JSON(http.StatusConflict, gin.H{"message": "User/Email already exists"})
		} else {
			log.Printf("Failed adding user: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Error"})
		}
		return
	}

	//TODO: Insert queue update here

	c.JSON(http.StatusCreated, user)
}

// UpdateUser creates a new user
func UpdateUser(c *gin.Context) {
	mc := c.MustGet("mongo").(mongo.Client)
	id := c.Param("id")

	var payload = &types.UserPatch{}
	errs := c.ShouldBindJSON(payload)
	if errs != nil {
		e := fmt.Sprintf("Invalid payload received: %s", errs)
		c.JSON(http.StatusBadRequest, e)
		return
	}

	user, err := mc.UpdateUser(id, *payload)
	if err != nil {
		if mc.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Not Found"})
		} else {
			log.Printf("Failed updating user: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Error"})
		}
		return
	}

	//TODO: Insert queue update here

	c.JSON(http.StatusOK, user)
}

// DeleteUser creates a new user
func DeleteUser(c *gin.Context) {
	mc := c.MustGet("mongo").(mongo.Client)
	id := c.Param("id")

	err := mc.DeleteUser(id)
	if err != nil {
		if mc.IsNotFound(err) {
			c.Status(http.StatusNoContent)
		} else {
			log.Printf("Failed deleting user: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Server Error"})
		}
		return
	}

	//TODO: Insert queue update here

	c.Status(http.StatusNoContent)
}
