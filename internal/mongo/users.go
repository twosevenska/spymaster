package mongo

import (
	"log"
	"regexp"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	uuid "github.com/satori/go.uuid"

	"types"
)

const usersCollection = "users"

// ListUsers lists the users for a given customer with a certain query
func (c Client) ListUsers(exactSearch map[string]interface{}, partialSearch map[string]string, perPage, pageNumber int) ([]types.User, int, error) {
	criteria := bson.M{}
	for field, val := range exactSearch {
		criteria[field] = val
	}
	for field, val := range partialSearch {
		criteria[field] = bson.RegEx{Pattern: regexp.QuoteMeta(val), Options: "i"}
	}

	log.Printf("Mongo: Trying to find a user that matches criteria: %+v", criteria)

	collection := c.Database.C(usersCollection)
	query := safeFind(collection, criteria).Sort("nickname")
	total, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	if perPage > 0 && pageNumber > 0 {
		query = query.Skip(perPage * (pageNumber - 1)).Limit(perPage)
	}

	r := []types.User{}
	query.All(&r)
	return r, total, err
}

// CreateUser creates a user for a given customer
func (c Client) CreateUser(payload types.UserPost) (user types.User, err error) {
	collection := c.Database.C(usersCollection)
	now := time.Now().UTC()

	user = types.User{
		ID:        uuid.NewV4(),
		Nickname:  payload.Nickname,
		Password:  payload.Password,
		Email:     payload.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if payload.FirstName != nil {
		user.FirstName = *payload.FirstName
	}
	if payload.LastName != nil {
		user.LastName = *payload.LastName
	}
	if payload.Country != nil {
		user.Country = *payload.Country
	}

	err = collection.Insert(user)
	return
}

// UpdateUser updates a user for a given customer
func (c Client) UpdateUser(userID string, payload types.UserPatch) (user types.User, err error) {
	collection := c.Database.C(usersCollection)

	u, err := uuid.FromString(userID)
	if err != nil {
		err = ErrInvalidUUID
		return
	}

	payload.UpdatedAt = time.Now().UTC()

	criteria := bson.M{"_id": u}
	change := mgo.Change{
		Update:    bson.M{"$set": payload},
		ReturnNew: true,
	}
	_, err = safeFind(collection, criteria).Apply(change, &user)
	return
}

// DeleteUser deletes a user for a given customer
func (c Client) DeleteUser(userID string) (err error) {
	collection := c.Database.C(usersCollection)

	u, err := uuid.FromString(userID)
	if err != nil {
		err = ErrInvalidUUID
		return
	}

	criteria := bson.M{"_id": u}
	change := mgo.Change{
		Remove: true,
	}
	_, err = safeFind(collection, criteria).Apply(change, nil)
	return err
}
