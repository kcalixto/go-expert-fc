package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Category struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *Category {
	return &Category{db: db}
}

func (c *Category) Create(name, description string) (Category, error) {
	category := Category{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
	}

	_, err := c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", category.ID, category.Name, category.Description)
	if err != nil {
		return Category{}, nil
	}

	return category, nil
}
