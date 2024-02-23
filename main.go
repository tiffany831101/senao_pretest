package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/tiffany831101/senao_pretest/database"
	"github.com/tiffany831101/senao_pretest/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type AccountCreationResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
}

type PasswordVerificationResponse struct {
	Success bool   `json:"success"`
	Reason  string `json:"reason,omitempty"`
}

var (
	db          *gorm.DB
	redisClient *redis.Client
)

const (
	maxAttempts     = 5
	waitDuration    = time.Minute
	rateLimitPrefix = "rate_limit:"
)

func init() {

	database.InitDB("admin", "pwd", "my_db", "my-db")
	if database.DB == nil {
		log.Fatal("Failed to initialize the database")
	}

	db = database.DB

	db.AutoMigrate(&User{})

	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "pwd",
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Connected to Redis successfully")
}

func createAccountHandler(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "reason": "Invalid JSON payload"})
		return
	}

	isPwdValid := utils.IsPasswordComplex(newUser.Password)
	if !isPwdValid {
		c.JSON(http.StatusBadRequest, AccountCreationResponse{Success: false, Reason: "Password must be 8-32 characters long and include at least 1 uppercase letter, 1 lowercase letter, and 1 number."})
	}

	var count int64
	db.Model(&User{}).Where("username = ?", newUser.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, AccountCreationResponse{Success: false, Reason: "Username already exists"})
		return
	}

	db.Create(&newUser)
	c.JSON(http.StatusOK, AccountCreationResponse{Success: true})
}

func verifyAccountAndPasswordHandler(c *gin.Context) {
	username := c.Param("username")
	var inputUser User
	if err := c.ShouldBindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, PasswordVerificationResponse{Success: false, Reason: "Invalid JSON payload"})
		return
	}

	attempts, err := incrementAttempts(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, PasswordVerificationResponse{Success: false, Reason: "Rate limit check failed"})
		return
	}

	if attempts > maxAttempts {
		c.JSON(http.StatusTooManyRequests, PasswordVerificationResponse{Success: false, Reason: "Too many attempts. Try again later."})
		return
	}

	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, PasswordVerificationResponse{Success: false, Reason: "Username not found"})
		return
	}

	if inputUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, PasswordVerificationResponse{Success: false, Reason: "Incorrect password"})
		return
	}

	err = resetAttempts(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, PasswordVerificationResponse{Success: false, Reason: "Failed to reset attempts"})
		return
	}

	c.JSON(http.StatusOK, PasswordVerificationResponse{Success: true})
}
func incrementAttempts(username string) (int, error) {
	key := rateLimitPrefix + username

	attempts, err := redisClient.Incr(context.Background(), key).Result()

	if err != nil {
		return 0, err
	}

	if attempts == 1 {
		redisClient.Expire(context.Background(), key, waitDuration)
	}

	return int(attempts), nil
}

func resetAttempts(username string) error {
	key := rateLimitPrefix + username
	_, err := redisClient.Del(context.Background(), key).Result()
	return err
}

func main() {
	router := gin.Default()

	router.POST("/account", createAccountHandler)

	router.POST("/account/:username/validate", verifyAccountAndPasswordHandler)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("Server is running on http://localhost%s\n", addr)
	router.Run(addr)

	defer database.CloseDB()
}
