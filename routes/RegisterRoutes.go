package routes

import "net/http"

func RegisterRoutes() {

	http.HandleFunc("GET /book", getAllBooks)
	http.HandleFunc("GET /book/{id}", getBookById)
	http.HandleFunc("POST /book", createBookById)
	http.HandleFunc("PUT /book/{id}", updateBook)
	http.HandleFunc("DELETE /book/{id}", deleteBook)

	http.HandleFunc("GET /", handleLivenessCheck)

}
