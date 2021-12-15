package server

import (
	"github.com/gin-gonic/gin"

	"internal/controllers"
	"internal/mongo"
)

// ContextParams holds the objects required
type ContextParams struct {
	MongoClient *mongo.Client
}

// CreateRouter creates a gin Engine and attaches backend clients, if needed as well as routes
// to all enpoints
func CreateRouter(mc *mongo.Client) *gin.Engine {
	contextParams := ContextParams{
		MongoClient: mc,
	}

	r := gin.Default()
	r.Use(ContextObjects(&contextParams))

	r.GET("/ping", controllers.Ping)
	r.GET("/health", controllers.Health)

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
