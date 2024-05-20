package main

import (
	"github.com/nullchefo/books-manager-project-go-lang/routes"
	"log"
	"net/http"
)

// why "gin" when we have the standard library?
// https://go.dev/blog/routing-enhancements
func main() {

	routes.RegisterRoutes()

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
