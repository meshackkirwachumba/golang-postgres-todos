package repository

import (
	"time"

	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/models"
)

func CreateUserInDB(pool *pgxpool.Pool, user *models.User) (*models.User, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  query := `
   INSERT INTO users_table (email, password)
   VALUES ($1, $2)
   RETURNING id, email, created_at, updated_at;
   `
    
   
   err :=pool.QueryRow(ctx, query, user.Email, user.Password).Scan(
	&user.ID,
	&user.Email,
	&user.CreatedAt,
	&user.UpdatedAt,
   )

   if err != nil {
	return nil, err
   }
  
   return user, nil
}

func GetUserByEmailInDB(pool *pgxpool.Pool, email string) (*models.User, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  query := `
   SELECT id, email, password, created_at, updated_at
   FROM users_table
   WHERE email = $1;
   `
	
   var user models.User

   err :=pool.QueryRow(ctx, query, email).Scan(
	&user.ID,
	&user.Email,
	&user.Password, 
	&user.CreatedAt,
	&user.UpdatedAt,
   )

   if err != nil {
	return nil, err
   }
  
   return &user, nil
}

func GetUserByIDInDB(pool *pgxpool.Pool, id int) (*models.User, error) {
  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()
  
  query := `
   SELECT id, email, password, created_at, updated_at
   FROM users_table
   WHERE id = $1;
   `
   var user models.User

   err :=pool.QueryRow(ctx, query, id).Scan(
	&user.ID,
	&user.Email,
	&user.Password, 
	&user.CreatedAt,
	&user.UpdatedAt,
   )

   if err != nil {
	return nil, err
   }
   
   return &user, nil
}