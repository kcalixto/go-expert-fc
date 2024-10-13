package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/kcalixto/go-expert-fc/sqlc/internal/db"
)

func main() {
	ctx := context.Background()

	dbConn, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/courses")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	queries := db.New(dbConn)

	err = queries.CreateCategory(ctx, db.CreateCategoryParams{
		ID:   uuid.New().String(),
		Name: gofakeit.Name(),
		Description: sql.NullString{
			String: "All about web development",
			Valid:  true,
		},
	})
	if err != nil {
		panic(err)
	}

	categories, err := queries.ListCategories(ctx)
	if err != nil {
		panic(err)
	}

	for _, category := range categories {
		fmt.Println(category.Name)
	}
}
