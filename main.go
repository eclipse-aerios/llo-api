package main

import (
	"log"

	"github.com/eclipse-aerios/llo-api/router"
)

func init() {
	log.Println("Running initialization function")
	// port := os.Getenv("PORT")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println("aeriOS LLO Operators API")
	log.Println("Developed in Go (GinGonic) and using K8s client for Go")

	app := router.NewRouter()
	app.Run(":8090")
}
