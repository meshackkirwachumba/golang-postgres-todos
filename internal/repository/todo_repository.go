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


func GetAllTodosFromDB(pool *pgxpool.Pool) ([]models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
      SELECT id, title, completed, created_at, updated_at
      FROM todos_table
      ORDER BY created_at DESC;
  `

  
  rows, err := pool.Query(ctx, query)
  if err != nil {
    return nil, err
  }
  defer rows.Close()

  var todoList []models.Todo

  for rows.Next() {
    var todo models.Todo
    err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
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

func GetTodoByIDFromDB(pool *pgxpool.Pool, id int) (*models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()  

  var query string = `
      SELECT id, title, completed, created_at, updated_at
      FROM todos_table
      WHERE id = $1;
  `

  var todo models.Todo
  err := pool.QueryRow(ctx, query, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
  if err != nil {
    return nil, err
  }
  return &todo, nil
}

func UpdatedTodoInDB(pool *pgxpool.Pool, id int, title string, completed bool) (*models.Todo, error) {
  var ctx context.Context 
  var cancel context.CancelFunc
  ctx, cancel = context.WithTimeout(context.Background(), 5 * time.Second)
  defer cancel()

  var query string = `
      UPDATE todos_table
      SET title = $1, completed = $2, updated_at = CURRENT_TIMESTAMP
      WHERE id = $3
      RETURNING id, title, completed, created_at, updated_at;
  `

  var todo models.Todo
  err := pool.QueryRow(ctx, query, title, completed, id).Scan(
    &todo.ID,
    &todo.Title,
    &todo.Completed,
    &todo.CreatedAt,
    &todo.UpdatedAt)
    
  if err != nil {
    return nil, err
  }

  return &todo, nil
}