package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/config"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/models"
	"github.com/meshackkirwachumba/golang-postgres-todos/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 6 characters long"})
			return
		}

	

		//hash the password before storing it in the database
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		user := &models.User{
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		createdUser, err:= repository.CreateUserInDB(pool, user)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
				c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return
		}



		
		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "user": createdUser})
	}
}

func LoginUserHandler(pool *pgxpool.Pool, config *config.ConfigStr) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserRequest LoginUserRequest
		if err := c.ShouldBindJSON(&loginUserRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"invalid json request", "details": err.Error()})
			return
		}

		user, err := repository.GetUserByEmailInDB(pool, loginUserRequest.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials!"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserRequest.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials!"})
			return
		}

		//map[string]interface{}]{"user_id": user.ID, "email": user.Email}
		claims := jwt.MapClaims{
			"user_id": user.ID,
			"email": user.Email,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		}


		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate signed token" + err.Error()})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
	}
}