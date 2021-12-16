package types

import (
	"time"

	"github.com/globalsign/mgo/bson"
)

// User stores bson query result
type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	FirstName string        `bson:"first_name" json:"first_name"`
	LastName  string        `bson:"last_name" json:"last_name"`
	Nickname  string        `bson:"nickname" json:"nickname"`
	Password  string        `bson:"password" json:"password"`
	Email     string        `bson:"email" json:"email"`
	Country   string        `bson:"country" json:"country"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

// UsersResult stores GetUsers response
type UsersResult struct {
	Users      []User `json:"objects"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	TotalCount int    `json:"total_count"`
}

// UserPost holds body for user creation request
type UserPost struct {
	// ID is set internaly
	FirstName *string `bson:"first_name" json:"first_name"`
	LastName  *string `bson:"last_name" json:"last_name"`
	Nickname  string  `bson:"nickname" json:"nickname" binding:"required"`
	Password  string  `bson:"password" json:"password" binding:"required"`
	Email     string  `bson:"email" json:"email" binding:"required"`
	Country   *string `bson:"country" json:"country"`
}

// UserPatch holds body for user update request
type UserPatch struct {
	FirstName *string   `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  *string   `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Nickname  *string   `bson:"nickname,omitempty" json:"nickname,omitempty"`
	Password  *string   `bson:"password,omitempty" json:"password,omitempty"`
	Email     *string   `bson:"email,omitempty" json:"email,omitempty"`
	Country   *string   `bson:"country,omitempty" json:"country,omitempty"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}
