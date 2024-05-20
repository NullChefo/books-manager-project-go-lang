package routes

import (
	"net/http"
)

func RegisterRoutes() {

	http.HandleFunc("GET /allBooks", getAllBooks)    // for admin purposes
	http.HandleFunc("GET /books", getBooksPaginated) // 50 books per page limit
	http.HandleFunc("GET /books/{id}", getBookById)
	http.HandleFunc("POST /books", createBook)
	http.HandleFunc("PUT /books/{id}", updateBook)
	http.HandleFunc("DELETE /books/{id}", deleteBook)

	http.HandleFunc("GET /", handleLivenessCheck)

}
