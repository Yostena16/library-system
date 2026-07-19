package controllers // ← matches the folder name

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"loan-service/internal/database"
	"loan-service/internal/models"
	"loan-service/internal/utils"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Register godoc
// @Summary      Register a new member
// @Description  Create a new library member account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  RegisterInput  true  "Registration details"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /auth/register [post]

// Register creates a new member account.
func Register(c *gin.Context) {
	var input RegisterInput

	// 1. Read + validate the JSON body
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Hash the password (so we never store it in plain text)
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not hash password"})
		return
	}

	// 3. Build the member
	member := models.Member{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	// 4. Save to the database
	if err := database.DB.Create(&member).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already in use"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "member registered successfully",
		"member":  member,
	})
}

// LoginInput is the JSON we expect when logging in.
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary      Log in a member
// @Description  Authenticate and receive a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  LoginInput  true  "Login credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /auth/login [post]
// Login checks credentials and returns a JWT token.

func Login(c *gin.Context) {
	var input LoginInput

	// 1. Read + validate the body
	if err := c.ShouldBindJSON(&input); err != nil { //bind the Input from our request to input
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Find the member by email
	var member models.Member
	if err := database.DB.Where("email = ?", input.Email).First(&member).Error; err != nil { //? is like placeholder go add input.email safely and First(&member) fetches the first result found in the database that matches the query and stores it in the member struct
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// 3. Compare the given password with the stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// 4. Generate a JWT token
	token, err := utils.GenerateToken(member.ID, member.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token, //hand the token to the client
	})
}
