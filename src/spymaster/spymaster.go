package spymaster

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"

	"spymaster/src/mongo"
	"spymaster/types"
)

var (
	// ErrDup indicates that an entry already exists
	ErrDup = errors.New("entry already exists")

	// ErrNotFound indicates that an entry is missing
	ErrNotFound = errors.New("entry not found")
)

func ListUsers(c *gin.Context, exact map[string]interface{}, partial map[string]string, perPage, pageNumber int) (response types.UsersResult, err error) {
	mc := c.MustGet("mongo").(mongo.Client)

	users, totalCount, err := mc.ListUsers(exact, partial, perPage, pageNumber)
	if err != nil {
		log.Printf("ListUsers: %s", err)
		return
	}

	response = types.UsersResult{
		Page:       pageNumber,
		PerPage:    perPage,
		TotalCount: totalCount,
		Users:      users,
	}

	return
}

// CreateUser creates a new user
func CreateUser(c *gin.Context, payload *types.UserPost) (user types.User, err error) {
	mc := c.MustGet("mongo").(mongo.Client)

	user, err = mc.CreateUser(*payload)
	if err != nil {
		if mc.IsDup(err) {
			err = ErrDup
		}
		log.Printf("Failed adding user: %s", err)
		return
	}

	//TODO: Insert queue update here

	return
}

// UpdateUser updates a user
func UpdateUser(c *gin.Context, id string, payload *types.UserPatch) (user types.User, err error) {
	mc := c.MustGet("mongo").(mongo.Client)

	user, err = mc.UpdateUser(id, *payload)
	if err != nil {
		if mc.IsNotFound(err) {
			err = ErrNotFound
		}
		log.Printf("Failed updating user: %s", err)
		return
	}

	//TODO: Insert queue update here

	return
}

// DeleteUser deletes a user
func DeleteUser(c *gin.Context, id string) error {
	mc := c.MustGet("mongo").(mongo.Client)

	err := mc.DeleteUser(id)
	if err != nil {
		if mc.IsNotFound(err) {
			err = ErrNotFound
		}
		log.Printf("Failed deleting user: %s", err)
		return err
	}

	//TODO: Insert queue update here

	return nil
}
