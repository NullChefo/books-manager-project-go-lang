package models

import "github.com/nullchefo/books-manager-project-go-lang/db"

func GenerateTablesFromModels() error {
	// TODO check if i can get all the models from the models package and dynamically create tables

	// TODO add your models here to create tables for them
	listOfModels := []any{Book{}}
	for _, model := range listOfModels {
		err := db.CreateTable(model)
		if err != nil {
			return err
		}
	}

	return nil
}
