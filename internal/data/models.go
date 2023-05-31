package data

import "database/sql"

type Models struct {
	Books BookModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Books: BookModel{DB: db},
	}
}
