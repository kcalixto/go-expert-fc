// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
)

type Category struct {
	ID          string
	Name        string
	Description sql.NullString
}

type Course struct {
	ID          string
	CategoryID  string
	Name        string
	Description sql.NullString
	Price       string
}
