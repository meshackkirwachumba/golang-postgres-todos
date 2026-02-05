package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/repository"
)

type CreateTodoInput struct {
	Title     string `json:"title" binding:"required"`
	Completed bool   `json:"completed"`
}

type UpdateTodoInput struct {
	Title     string `json:"title"`
	Completed *bool   `json:"completed"`
}

func CreateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
   return func(c *gin.Context) {
	  var input CreateTodoInput
	  if err := c.ShouldBindJSON(&input); err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input json"})
		 return
	  }

	  todo, err := repository.CreateTodoInDB(pool, input.Title, input.Completed)
	  if err != nil {
		 c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save todo in database"})
		 return
	  }

	  c.JSON(http.StatusCreated, todo)
   }
}


func GetAllTodosHandler(pool *pgxpool.Pool) gin.HandlerFunc {
  return func(c *gin.Context) {
	todos, err := repository.GetAllTodosFromDB(pool)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todos from database"})
		return
	}
	c.JSON(http.StatusOK, todos)
  }
}

func GetTodoByIDHandler(pool *pgxpool.Pool) gin.HandlerFunc {
  return func(c *gin.Context) {
	idString := c.Param("id")

	todoID, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	todo, err := repository.GetTodoByIDFromDB(pool, todoID)
	if err != nil {
		if err == pgx.ErrNoRows{
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve todo from database"})
		return
	}

	c.JSON(http.StatusOK, todo)
  }
}

func UpdateTodoHandler(pool *pgxpool.Pool) gin.HandlerFunc {
  return func(c *gin.Context) {
	idString := c.Param("id")

	todoID, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid todo ID"})
		return
	}

	var input UpdateTodoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input json"})
		return
	}

	if input.Title == "" && input.Completed == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one field (title or completed) must be provided for update"})
		return
	}

	var completed bool
	if input.Completed != nil {
		completed = *input.Completed
	}

	todo, err := repository.UpdateTodoInDB(pool, todoID, input.Title, completed)
	if err != nil {
		if err == pgx.ErrNoRows{
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo to be updated not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update todo in database"})
		return
	}

	c.JSON(http.StatusOK, todo)
  }
}