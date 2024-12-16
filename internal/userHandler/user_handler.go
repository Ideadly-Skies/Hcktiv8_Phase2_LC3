package handler

import (
	"fmt"
	"net/http"
	config "w4/lc3/config/database"
	
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgconn"
)

// -- Create Users table
// CREATE TABLE Users (
//     user_id SERIAL PRIMARY KEY,
//     name VARCHAR(100),
//     email VARCHAR(100) UNIQUE NOT NULL,
//     password VARCHAR(255) NOT NULL,
//     jwt_token VARCHAR(255)
// );

// Users struct
type Users struct {
	ID  int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	JwtToken string `json:"jwt_token"`
}

// RegisterRequest struct
type RegisterRequest struct {
	Name string `json:"name" validate:"required,name"` 				// Name of the user
	Email    string `json:"email" validate:"required,email"`        // Email address
	Password string `json:"password" validate:"required,password"`  // Password for the account
}

// login request struct
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// login response: token
type LoginResponse struct {
	Token string `json:"token"`
}

var jwtSecret = []byte("12345")

// @Summary Register a new user
// @Description Create a new user account by providing name, email, and password
// @Tags Users
// @Accept  json
// @Produce  json
// @Param request body RegisterRequest true "User registration data"
// @Success 201 {object} map[string]interface{} "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid input or email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/register [post]
func Register(c echo.Context) error {
    var req RegisterRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Request"})
    }

	// hash the password
    hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
    }

    // queries to insert to both users and customers db
	users_query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING user_id"
	
	var userID int
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// query row 1: insert to users 
	err = config.Pool.QueryRow(ctx, users_query, req.Name, req.Email, string(hashPassword)).Scan(&userID)
	if err != nil {
		fmt.Println("Error inserting into users table:", err)

		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" { // Unique violation (email already registered)
				return c.JSON(http.StatusBadRequest, map[string]string{"message": "Email already registered"})
			}
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Internal Server Error"})
	}

    return c.JSON(http.StatusOK, map[string]interface{}{
        "message": "User registered successfully",
        "user_id": string(userID),
        "email": req.Email,
    })
}

// @Summary Login a user
// @Description Authenticate a user by providing valid credentials
// @Tags Users
// @Accept  json
// @Produce  json
// @Param request body LoginRequest true "User login data"
// @Success 200 {object} LoginResponse "Authentication successful with a JWT token"
// @Failure 400 {object} map[string]string "Invalid email or password"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/login [post]
func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message":"Invalid Request"})
	}
	
	var user Users
	query := "SELECT user_id, email, password FROM users WHERE email = $1"
	err := config.Pool.QueryRow(context.Background(), query, req.Email).Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email or password"})
	}

	// compare password to see if it matches the student password provided
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid email or password"})
	}

	// create new jwt claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // Use `jwt.NewNumericDate` for expiry
	})
	
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Invalid Generate Token"})
	}

	// Update the jwt_token column in the database
	updateQuery := "UPDATE users SET jwt_token = $1 WHERE user_id = $2"
	_, err = config.Pool.Exec(context.Background(), updateQuery, tokenString, user.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update token"})
	}

	// return ok status and login response
	return c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}