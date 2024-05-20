package routes

import (
	"encoding/json"
	"fmt"
	"github.com/nullchefo/books-manager-project-go-lang/models"
	"github.com/nullchefo/books-manager-project-go-lang/utils"
	"net/http"
	"strings"
)

func getAllBooks(w http.ResponseWriter, r *http.Request) {
	if !utils.CheckAdmin(w, r) {
		return
	}

	books, err := models.GetAllBooks()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, books, http.StatusOK)
}

func getBooksPaginated(w http.ResponseWriter, r *http.Request) {
	params := utils.ParsePaginationParams(r)
	offset := (params.Page - 1) * params.Size

	if params.Size > 50 {
		http.Error(w, "Limit must be less than 50", http.StatusBadRequest)
		return
	}

	books, err := utils.GetPaginatedData(offset, params.Size, models.GetBooksPaginated)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	totalBooks, err := models.GetTotalBooksCount()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to retrieve books count", http.StatusInternalServerError)
		return
	}

	response := utils.CreatePaginatedResponse(books, params.Page, params.Size, totalBooks)
	utils.SendJSONResponse(w, response, http.StatusOK)
}

func getBookById(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := models.GetBookById(id)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "no rows in result set") {

			http.Error(w, "Book not found", http.StatusNotFound)
		} else {

			http.Error(w, "Failed to retrieve book", http.StatusInternalServerError)
		}
		return
	}

	utils.SendJSONResponse(w, book, http.StatusOK)
}

func createBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Println(book)

	err = book.Save()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	utils.SendJSONResponse(w, book, http.StatusCreated)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	book.ID = id
	updatedBook, err := models.UpdateBookById(id, book)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update book", http.StatusInternalServerError)
		}
		return
	}

	utils.SendJSONResponse(w, updatedBook, http.StatusOK)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(r)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteBookById(id)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Book not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK) // OK
}
