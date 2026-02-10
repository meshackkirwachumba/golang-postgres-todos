package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/models"
)

func CreateTodoInDB(pool *pgxpool.Pool, title string, completed bool, userID string) (*models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
    INSERT INTO todos_table (title, completed, user_id)
	  VALUES ($1, $2, $3)
	  RETURNING id, title, completed, created_at, updated_at, user_id;
  `

  var todoRow models.Todo

  err :=pool.QueryRow(ctx, query, title, completed, userID).Scan(
    &todoRow.ID,
    &todoRow.Title,
    &todoRow.Completed,
    &todoRow.CreatedAt,
    &todoRow.UpdatedAt,
    &todoRow.UserID,
  )

  if err != nil {
    return nil, err
  }

  return &todoRow, nil
}	


func GetAllTodosFromDB(pool *pgxpool.Pool, userID string) ([]models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
      SELECT id, title, completed, created_at, updated_at, user_id
      FROM todos_table
      WHERE user_id = $1
      ORDER BY created_at DESC;
  `

  
  rows, err := pool.Query(ctx, query, userID)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var todoList []models.Todo

  for rows.Next() {
    var todo models.Todo
    err := rows.Scan(
      &todo.ID,
      &todo.Title,
      &todo.Completed,
      &todo.CreatedAt,
      &todo.UpdatedAt,
      &todo.UserID,
    )
    if err != nil {
      return nil, err
    }
    todoList = append(todoList, todo)
  }

  if err = rows.Err(); err != nil {
    return nil, err
  }

  return todoList, nil
}

func GetTodoByIDFromDB(pool *pgxpool.Pool, id int, userID string) (*models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()  

  var query string = `
      SELECT id, title, completed, created_at, updated_at, user_id
      FROM todos_table
      WHERE id = $1 AND user_id = $2;
  `

  var todo models.Todo
  err := pool.QueryRow(ctx, query, id, userID).Scan(
    &todo.ID,
    &todo.Title,
    &todo.Completed,
    &todo.CreatedAt,
    &todo.UpdatedAt,
    &todo.UserID,
  )
  if err != nil {
    return nil, err
  }
  return &todo, nil
}

func UpdateTodoInDB(pool *pgxpool.Pool, id int, title string, completed bool, userID string) (*models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
      UPDATE todos_table
      SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
      WHERE id = $3 AND user_id = $4
      RETURNING id, title, completed, created_at, updated_at, user_id;
  `

  var todo models.Todo
  err := pool.QueryRow(ctx, query, title, completed, id, userID).Scan(
    &todo.ID,
    &todo.Title,
    &todo.Completed,
    &todo.CreatedAt,
    &todo.UpdatedAt,
    &todo.UserID)

  if err != nil {
    return nil, err
  }

  return &todo, nil
}

func DeleteTodoFromDB(pool *pgxpool.Pool, id int, userID string) error {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
      DELETE FROM todos_table
      WHERE id = $1 AND user_id = $2;
  `

  commandTag, err := pool.Exec(ctx, query, id, userID)

  if err != nil {
    return err
  }

  if commandTag.RowsAffected() == 0 {
    return  fmt.Errorf("todo with id %d not found", id)
  } 
  
  return nil

}
