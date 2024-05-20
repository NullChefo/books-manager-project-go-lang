package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"reflect"
	"strings"
)

var db *sql.DB

func GetDb() *sql.DB {
	return db
}

func Init(persistData bool) error {

	if persistData == false {
		err := os.Remove("example.db")
		return err
	}

	var err error
	db, err = sql.Open("sqlite", "example.db")
	if err != nil {

		fmt.Println(err)
		fmt.Println("===============")
		panic("No database connection")
	}
	db.SetMaxOpenConns(10)

	return nil
}

// Creates a new table based on the provided structure
func CreateTable(dataType any) error {
	// get the data for the type
	typeRef := reflect.TypeOf(dataType)

	var tableColumns []string

	for i := 0; i < typeRef.NumField(); i++ {
		field := typeRef.Field(i)

		if field.Tag.Get("db") == "transient" {
			continue
		}

		columnName := strings.ToLower(field.Name)
		columnType, err := getSQLType(field.Type)
		if err != nil {
			return err
		}

		// Add Primary Key and Auto Increment
		if field.Tag.Get("db") == "primaryKey" {
			tableColumns = append(tableColumns, fmt.Sprintf("%s %s PRIMARY KEY AUTOINCREMENT", columnName, columnType))
			continue
		}

		tableColumns = append(tableColumns, fmt.Sprintf("%s %s", columnName, columnType))
	}

	tableName := strings.ToLower(typeRef.Name() + "s")
	queryString := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)", tableName, strings.Join(tableColumns, ", "))

	fmt.Println(queryString)

	_, err := db.Exec(queryString)
	if err != nil {
		return err
	}

	return nil
}

// Creates a new INSERT query based on the structure and returns an object
// Insert creates a new INSERT query based on the structure and returns the provided object with the ID set.
func Insert[T any](dataType T) (T, error) {
	// get the data for the type
	typeRef := reflect.TypeOf(dataType)
	valueRef := reflect.ValueOf(dataType)

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

	fmt.Println(queryString)

	statement, err := db.Prepare(queryString)
	if err != nil {
		return dataType, err
	}

	defer statement.Close()

	res, err := statement.Exec(statementValue...)
	if err != nil {
		return dataType, err
	}

	_, err = res.LastInsertId()
	if err != nil {
		return dataType, err
	}

	// set the primary key of the object

	return dataType, nil
}

func getSQLType(fieldType reflect.Type) (string, error) {
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "INTEGER", nil
	case reflect.Float32, reflect.Float64:
		return "REAL", nil
	case reflect.Bool:
		return "BOOLEAN", nil
	case reflect.String:
		return "TEXT", nil
	default:
		return "", fmt.Errorf("unsupported field type: %s", fieldType.Kind())
	}

}

// Generic Update function to update a record in the database
func Update[T any](dataType T) (T, error) {
	var updatedObject T

	typeRef := reflect.TypeOf(dataType)
	valueRef := reflect.ValueOf(dataType)

	var tableColumns []string
	var statementValues []any
	var primaryKeyField string
	var primaryKeyValue any

	for i := 0; i < typeRef.NumField(); i++ {
		field := typeRef.Field(i)

		if field.Tag.Get("db") == "primaryKey" {
			primaryKeyField = strings.ToLower(field.Name)
			primaryKeyValue = valueRef.Field(i).Interface()
			continue
		}

		if field.Tag.Get("db") == "transient" {
			continue
		}

		tableColumns = append(tableColumns, strings.ToLower(field.Name)+" = ?")
		statementValues = append(statementValues, valueRef.Field(i).Interface())
	}

	if primaryKeyField == "" {
		return updatedObject, fmt.Errorf("primary key field not found")
	}

	statementValues = append(statementValues, primaryKeyValue)

	tableName := strings.ToLower(typeRef.Name() + "s")
	setClause := strings.Join(tableColumns, ", ")
	queryString := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?", tableName, setClause, primaryKeyField)

	statement, err := db.Prepare(queryString)
	if err != nil {
		return updatedObject, err
	}
	defer statement.Close()

	_, err = statement.Exec(statementValues...)
	if err != nil {
		return updatedObject, err
	}

	// Fetch the updated object from the database
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", tableName, primaryKeyField)
	row := db.QueryRow(selectQuery, primaryKeyValue)

	columnValues := make([]interface{}, typeRef.NumField())
	for i := 0; i < typeRef.NumField(); i++ {
		columnValues[i] = valueRef.Field(i).Addr().Interface()
	}

	err = row.Scan(columnValues...)
	if err != nil {
		return updatedObject, err
	}

	return dataType, nil
}

// Generics baabyyyyy
// Creates a new SELECT query based on the structure and returns an array of the same structure
func Select[T any](dataType T) ([]T, error) {
	// Create a slice to hold the results
	var results []T

	// get the data for the type
	typeRef := reflect.TypeOf(dataType)

	tableName := strings.ToLower(typeRef.Name() + "s")

	queryString := fmt.Sprintf("SELECT * FROM %s", tableName)

	dbCursor, err := db.Query(queryString)
	if err != nil {
		return nil, err
	}

	for dbCursor.Next() {
		// Create a new object to hold the row data
		newObject := reflect.New(typeRef).Elem()

		// Create a slice of interfaces to hold the column values
		columnValues := make([]interface{}, typeRef.NumField())
		for i := 0; i < typeRef.NumField(); i++ {
			columnValues[i] = newObject.Field(i).Addr().Interface()
		}

		err = dbCursor.Scan(columnValues...)
		if err != nil {
			return nil, err
		}

		results = append(results, newObject.Interface().(T))

	}

	return results, nil

}

func Delete[T any](dataType T, id int64) error {

	// Get the type information of the provided data type
	typeRef := reflect.TypeOf(dataType)

	var primaryKeyField string

	for i := 0; i < typeRef.NumField(); i++ {
		field := typeRef.Field(i)

		// Identify the primary key field
		if field.Tag.Get("db") == "primaryKey" {
			primaryKeyField = strings.ToLower(field.Name)
			continue
		}
	}

	if primaryKeyField == "" {
		return fmt.Errorf("primary key field not found")
	}

	// Construct the table name from the type name
	tableName := strings.ToLower(typeRef.Name() + "s")
	queryString := fmt.Sprintf("DELETE FROM %s WHERE %s = %d", tableName, primaryKeyField, id)
	fmt.Println(queryString)

	statement, err := db.Prepare(queryString)
	if err != nil {
		return err
	}

	_, err = statement.Exec(queryString)
	if err != nil {
		return err
	}

	defer statement.Close()

	return nil
}

func SelectById[T any](dataType T, id int64) (T, error) {

	// Get the type information of the provided data type
	typeRef := reflect.TypeOf(dataType)

	var primaryKeyField string

	for i := 0; i < typeRef.NumField(); i++ {
		field := typeRef.Field(i)

		// Identify the primary key field
		if field.Tag.Get("db") == "primaryKey" {
			primaryKeyField = strings.ToLower(field.Name)
			continue
		}
	}

	if primaryKeyField == "" {
		return *new(T), fmt.Errorf("primary key field not found")
	}

	// Construct the table name from the type name
	tableName := strings.ToLower(typeRef.Name() + "s")

	// Construct the SQL query string
	queryString := fmt.Sprintf("SELECT * FROM %s WHERE %s = %d", tableName, primaryKeyField, id)

	// Execute the query
	dbCursor, err := db.Query(queryString)
	if err != nil {
		return *new(T), err
	}
	defer dbCursor.Close()

	// Check if a result was found
	if !dbCursor.Next() {
		return *new(T), sql.ErrNoRows
	}

	// Create a new instance of the provided type
	newObject := reflect.New(typeRef).Elem()

	// Create a slice of interfaces to hold the column values
	columnValues := make([]interface{}, typeRef.NumField())
	for i := 0; i < typeRef.NumField(); i++ {
		columnValues[i] = newObject.Field(i).Addr().Interface()
	}

	// Scan the row into the new object
	if err := dbCursor.Scan(columnValues...); err != nil {
		return *new(T), err
	}

	return newObject.Interface().(T), nil
}

// GetTotalCount creates a new SELECT COUNT query based on the structure and returns the total count of records.
func GetTotalCount[T any](dataType T) (int, error) {
	// Get the data for the type
	typeRef := reflect.TypeOf(dataType)

	tableName := strings.ToLower(typeRef.Name() + "s")

	queryString := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)

	var count int
	err := db.QueryRow(queryString).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SelectPaginated creates a new SELECT query based on the structure and returns a paginated array of the same structure.
func SelectPaginated[T any](dataType T, offset int, limit int) ([]T, error) {
	// Create a slice to hold the results
	var results []T

	// Get the data for the type
	typeRef := reflect.TypeOf(dataType)

	tableName := strings.ToLower(typeRef.Name() + "s")

	queryString := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", tableName)

	dbCursor, err := db.Query(queryString, limit, offset)
	if err != nil {
		return nil, err
	}
	defer dbCursor.Close()

	for dbCursor.Next() {
		// Create a new object to hold the row data
		newObject := reflect.New(typeRef).Elem()

		// Create a slice of interfaces to hold the column values
		columnValues := make([]interface{}, typeRef.NumField())
		for i := 0; i < typeRef.NumField(); i++ {
			columnValues[i] = newObject.Field(i).Addr().Interface()
		}

		err = dbCursor.Scan(columnValues...)
		if err != nil {
			return nil, err
		}

		results = append(results, newObject.Interface().(T))
	}

	return results, nil
}
