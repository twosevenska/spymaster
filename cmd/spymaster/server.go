package main

import (
	"github.com/gin-gonic/gin"

	"internal/controllers"
)

// CreateRouter creates a gin Engine and attaches backend clients, if needed as well as routes
// to all enpoints
func createRouter(contextParams *ContextParams) *gin.Engine {
	r := gin.Default()
	r.Use(ContextObjects(contextParams))

	r.GET("/ping", controllers.Ping)

	api := r.Group("/", controllers.Pagination())
	{
		api.GET("/users", controllers.ListUsers)
		api.POST("/user", controllers.CreateUser)
		api.PATCH("/user", controllers.UpdateUser)
		api.DELETE("/user", controllers.DeleteUser)
	}

	return r
}

// ContextObjects attaches backend clients to the API context
func ContextObjects(contextParams *ContextParams) gin.HandlerFunc {
	return func(c *gin.Context) {
		newMongo := contextParams.MongoClient.Copy()
		defer newMongo.Close()
		c.Set("mongo", newMongo)
		c.Next()
	}
}
