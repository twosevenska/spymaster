package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"

	"internal/mongo"
)

// Config ...
type Config struct {
	// BaseURL   string       `envconfig:"base_url" default:"https://10.0.0.1"`
	// UserTopic string       `envconfig:"user_notification_topic" default:"user-notifications"`
	Mongo mongo.Config `envconfig:"mongo"`
}

// ContextParams holds the objects required
type ContextParams struct {
	MongoClient *mongo.Client
}

func main() {
	var conf Config
	err := envconfig.Process("spymaster", &conf)
	if err != nil {
		log.Fatalf("Failed to load env config: %s", err.Error())
	}

	Splash()

	mc, err := mongo.Connect(conf.Mongo)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	contextParams := ContextParams{
		MongoClient: mc,
	}

	r := createRouter(&contextParams)
	r.Run(":7000") // listen and serve on 0.0.0.0:7000
}
