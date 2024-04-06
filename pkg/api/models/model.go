package models

import (
	"database/sql"
	"errors"
	"log"
	"os"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	User     UserModel
	Order    OrderModel
	Product  ProductModel
	Category CategoryModel
}

func NewModels(db *sql.DB) Models {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	return Models{
		User: UserModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Order: OrderModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Product: ProductModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
		Category: CategoryModel{
			DB:       db,
			InfoLog:  infoLog,
			ErrorLog: errorLog,
		},
	}
}
