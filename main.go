package main

import (
	"flag"
	"fmt"
	"net/http"

	"bitbucket.org/hebertthome/traning-oauth-go/config"
	"bitbucket.org/hebertthome/traning-oauth-go/context"
	"bitbucket.org/hebertthome/traning-oauth-go/handlers"
	"bitbucket.org/hebertthome/traning-oauth-go/logger"
	"bitbucket.org/hebertthome/traning-oauth-go/redis"
)

func main() {
	fmt.Printf("%v\n", "let's go to begin ... \n\n")

	// Load Configuration by yaml file
	var path string
	flag.StringVar(&path, "conf", "", "yaml configuration file")
	flag.Parse()

	if err := config.Setup(path); err != nil {
		panic(err)
	}
	configuration := config.Get()

	// Start Redis
	cache := redis.Start(&configuration)

	// Create APP Context
	ctx := &context.AppContext{Cache: cache, Logger: logger.GetLogger()}

	// Configure Handlers
	http.Handle("/authenticate", handlers.AppHandler{C: ctx, H: handlers.Authenticate})
	http.Handle("/api", handlers.AppHandler{C: ctx, H: handlers.API, Auth: true})

	err := http.ListenAndServe(configuration.Bind, nil)
	if err != nil {
		panic(err)
	}

}
