package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kcalixto/go-expert-fc/sqlc/internal/db"
)

type CourseDB struct {
	db *sql.DB
	*db.Queries
}

func NewCourseDB(dbConn *sql.DB) *CourseDB {
	return &CourseDB{
		db:      dbConn,
		Queries: db.New(dbConn),
	}
}

func (c *CourseDB) CallTX(ctx context.Context, fn func(*db.Queries) error) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = fn(db.New(tx))
	if err != nil {
		if errRB := tx.Rollback(); errRB != nil {
			return fmt.Errorf("fn err: %v, rollback err: %v", err, errRB)
		}
		return err
	}
	return tx.Commit()
}

func main() {
	// ctx := context.Background()

	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	courseDB := NewCourseDB(dbConn)

	err = courseDB.CallTX(context.Background(), func(q *db.Queries) error {
		return q.CreateCategory(context.Background(), db.CreateCategoryParams{
			ID:          "1",
			Name:        "Web Development",
			Description: sql.NullString{String: "All about web development", Valid: true},
		})
	})
	if err != nil {
		panic(err)
	}

}
