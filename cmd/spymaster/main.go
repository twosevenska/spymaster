package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"

	"internal/mongo"
	"internal/server"
)

// Config ...
type Config struct {
	// UserTopic string       `envconfig:"user_notification_topic" default:"user-notifications"`
	Mongo mongo.Config `envconfig:"mongo"`
}

func main() {
	var conf Config
	err := envconfig.Process("spymaster", &conf)
	if err != nil {
		log.Fatalf("Failed to load env config: %s", err.Error())
	}

	fmt.Println(splash)

	mc, err := mongo.Connect(conf.Mongo)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %s", err)
	}

	r := server.CreateRouter(mc)
	r.Run(":7000") // listen and serve on 0.0.0.0:7000
}

const splash = `

  *****************************************
  *               Spymaster               *
  *****************************************
                ████████████              
              ██▓▓▓▓▓▓▓▓▓▓▓▓██            
            ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██          
            ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██        
          ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██        
          ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██        
          ██▓▓▓▓▓▓    ████   ▓▓██        
            ██▓▓▓▓    ████   ▓▓█        
          ████▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██        
      ████▓▓▓▓██▓▓  ████████  ██████      
    ██▓▓▓▓▓▓▓▓▓▓██          ██▓▓▓▓▓▓██    
    ██▓▓▓▓▓▓▓▓▓▓▓▓██████████▓▓▓▓▓▓▓▓██    
  ██▓▓▓▓▓▓▓▓██▓▓▓▓  ██  ▓▓▓▓██▓▓▓▓▓▓▓▓██  
  ██▓▓▓▓██████▓▓▓▓▓▓██  ▓▓▓▓██████▓▓▓▓██  
  ██▓▓▓▓▓▓████▓▓▓▓▓▓▓▓▓▓▓▓▓▓████▓▓▓▓▓▓██  
  ██▓▓▓▓▓▓████▓▓▓▓▓▓▓▓▓▓▓▓▓▓████▓▓▓▓▓▓██  
    ██████  ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓██  ██████    
          ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓██          
        ██▓▓▓▓▓▓▓▓▓▓██▓▓▓▓▓▓▓▓▓▓██        
      ████▓▓▓▓▓▓▓▓██  ██▓▓▓▓▓▓▓▓████      
  ████▓▓▓▓▓▓▓▓▓▓██      ██▓▓▓▓▓▓▓▓▓▓████  
██▓▓▓▓▓▓▓▓▓▓▓▓▓▓██      ██▓▓▓▓▓▓▓▓▓▓▓▓▓▓██
██████████████████      ██████████████████

`
