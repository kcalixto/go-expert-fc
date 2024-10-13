package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type Course struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) *Course {
	return &Course{db: db}
}

func (c *Course) Create(name, description, categoryID string) (Course, error) {
	course := Course{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		CategoryID:  categoryID,
	}

	_, err := c.db.Exec("INSERT INTO courses (id, name, description, category_id) VALUES ($1, $2, $3, $4)", course.ID, course.Name, course.Description, course.CategoryID)
	if err != nil {
		return Course{}, nil
	}

	return course, nil
}

func (c *Course) FindAll() ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	categories := []Course{}
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.CategoryID); err != nil {
			return nil, err
		}
		categories = append(categories, course)
	}

	return categories, nil
}

func (c *Course) FindByCategoryID(categoryID string) ([]Course, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, nil
	}
	defer rows.Close()
	categories := []Course{}
	for rows.Next() {
		var course Course
		if err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.CategoryID); err != nil {
			return nil, err
		}
		categories = append(categories, course)
	}

	return categories, nil
}
