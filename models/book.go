package models

import (
	"fmt"
	"github.com/nullchefo/books-manager-project-go-lang/db"
	"reflect"
	"strings"
)

type Book struct {
	ID                int64  `json:"id" db:"primaryKey"`
	Title             string `json:"title"`
	ISBN              string `json:"isbn"`
	Author            string `json:"author"`
	YearOfPublication int    `json:"year_of_publication"`
	//TestField         string `json:"test_field" db:"transient"`
}

func GetAllBooks() (books []Book, err error) {
	books, err = db.Select(Book{})
	if err != nil {
		return nil, err
	}
	return books, nil
}

func GetBookById(id int64) (book Book, err error) {
	book, err = db.SelectById(Book{}, id)
	return book, err
}

// UpdateBookById updates the book with the given ID using the provided book details.
// It returns the updated book and an error if the book is not found or if the update operation fails.
func UpdateBookById(id int64, book Book) (Book, error) {
	// Find the book by ID
	foundBook, err := db.SelectById(Book{}, id)
	if err != nil {
		return Book{}, fmt.Errorf("failed to find book by id: %w", err)
	}
	// Check if the book was found
	if foundBook.ID == 0 {
		return Book{}, fmt.Errorf("book not found")
	}

	// Ensure the book to update has the correct ID
	book.ID = id

	// Update the book using the generic Update function
	updatedBook, err := db.Update(book)
	if err != nil {
		return Book{}, fmt.Errorf("failed to update book: %w", err)
	}

	return updatedBook, nil
}

func DeleteBookById(id int64) (err error) {
	err = db.Delete(Book{}, id)

	if err != nil {
		return err
	}

	return nil
}

func GetBooksPaginated(page int, pageSize int) (books []Book, err error) {
	books, err = db.SelectPaginated(Book{}, page, pageSize)
	if err != nil {
		return nil, err
	}
	return books, nil

}

func GetTotalBooksCount() (int, error) {
	return db.GetTotalCount(Book{})
}

// TODO make this happen
//func SaveBook(book Book) (Book, error) {
//	res, err := db.Insert(book)
//	if err != nil {
//		return Book{}, err
//	}
//	return res, nil
//
//}

func (book *Book) Save() error {

	// get the data for the type
	typeRef := reflect.TypeOf(*book)   // WTF
	valueRef := reflect.ValueOf(*book) // WTF

	var tableColumns []string
	var statementValue []any
	var tableValues []string

	for i := 0; i < typeRef.NumField(); i++ {
		if typeRef.Field(i).Tag.Get("db") == "transient" {
			continue
		}

		if typeRef.Field(i).Tag.Get("db") == "primaryKey" {
			continue
		}

		tableColumns = append(tableColumns, strings.ToLower(typeRef.Field(i).Name))
		tableValues = append(tableValues, "?")
		statementValue = append(statementValue, valueRef.Field(i).Interface())
	}

	tableName := strings.ToLower(typeRef.Name() + "s")
	queryStringColumn := "(" + strings.Join(tableColumns, ",") + ")"
	queryStringValues := "(" + strings.Join(tableValues, ",") + ")"

	queryString := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, queryStringColumn, queryStringValues)

	statement, err := db.GetDb().Prepare(queryString)

	defer statement.Close()

	if err != nil {
		return err
	}

	result, err := statement.Exec(statementValue...)
	if err != nil {
		return err
	}

	bookId, err := result.LastInsertId()
	book.ID = bookId

	return err
}
