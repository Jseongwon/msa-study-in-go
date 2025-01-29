package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"noti-study/configs"

	"firebase.google.com/go/messaging"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=username dbname=mydb sslmode=disable password=mypassword"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func registerDeviceToken(c *gin.Context) {
	email := c.PostForm("email")
	token := c.PostForm("token")

	_, err := db.Exec("INSERT INTO device_tokens (email, token) VALUES ($1, $2) ON CONFLICT (email) DO UPDATE SET token = $2", email, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device token registered successfully"})
}

func sendPushNotification(c *gin.Context) {
	email := c.PostForm("email")
	targetEmail := c.PostForm("target_email")

	var targetToken string
	err := db.QueryRow("SELECT token FROM device_tokens WHERE email = $1", targetEmail).Scan(&targetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve target token"})
		return
	}

	fcmConfig := configs.LoadFCMConfig()
	opt := option.WithAPIKey(fcmConfig.APIKey)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v", err)
	}

	message := &messaging.Message{
		Token: targetToken,
		Notification: &messaging.Notification{
			Title: "New Message",
			Body:  fmt.Sprintf("You have a new message from %s", email),
		},
	}

	_, err = client.Send(context.Background(), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send push notification"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Push notification sent successfully"})
}

func main() {
	initDB()

	r := gin.Default()

	r.POST("/register", registerDeviceToken)
	r.POST("/send", sendPushNotification)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Notification server running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
