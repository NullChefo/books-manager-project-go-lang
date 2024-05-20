package main

import (
	"github.com/nullchefo/books-manager-project-go-lang/db"
	"github.com/nullchefo/books-manager-project-go-lang/models"
	"github.com/nullchefo/books-manager-project-go-lang/routes"
	"log"
	"net/http"
)

// why "gin" when we have the standard library?
// https://go.dev/blog/routing-enhancements
func main() {
	// Initialize the database
	err := db.Init(true)
	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	// Create the tables
	err = models.GenerateTablesFromModels()
	if err != nil {
		log.Fatal("Could not create the tables")
	}

	// Register the routes
	routes.RegisterRoutes()
	// Start the server
	log.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
		return
	}
}
