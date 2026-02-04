package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/models"
)

func CreateTodoInDB(pool *pgxpool.Pool, title string, completed bool) (*models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
      INSERT INTO todos_table (title, completed)
	  VALUES ($1, $2)
	  RETURNING id, title, completed, created_at, updated_at;
  `

  var todoRow models.Todo

  err :=pool.QueryRow(ctx, query, title, completed).Scan(&todoRow.ID,&todoRow.Title,&todoRow.Completed,&todoRow.CreatedAt,&todoRow.UpdatedAt,
  )

  if err != nil {
    return nil, err
  }

  return &todoRow, nil
}	