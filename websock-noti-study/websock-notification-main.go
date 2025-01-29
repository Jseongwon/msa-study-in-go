package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"redis-study/configs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

var db *sql.DB
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
	Body   string `json:"body"`
}

func initDB() {
	dbConfig := configs.LoadDBConfig()
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		dbConfig.User, dbConfig.DBName, dbConfig.Password, dbConfig.Host, dbConfig.Port)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

func registerUser(c *gin.Context) {
	email := c.PostForm("email")
	userID := uuid.New().String()

	_, err := db.Exec("INSERT INTO users (id, email) VALUES ($1, $2) ON CONFLICT (email) DO NOTHING", userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "user_id": userID})
}

func registerDeviceToken(c *gin.Context) {
	userID := c.PostForm("user_id")
	token := c.PostForm("token")

	_, err := db.Exec("INSERT INTO device_tokens (user_id, token) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET token = $2", userID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register device token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device token registered successfully"})
}

func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, ws)
			break
		}

		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		// Store message in the database
		_, err := db.Exec("INSERT INTO messages (id, room_id, user_id, body) VALUES ($1, $2, $3, $4)",
			uuid.New().String(), msg.RoomID, msg.UserID, msg.Body)
		if err != nil {
			log.Println("Failed to store message:", err)
			continue
		}

		// Send message to all connected clients in the room
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}

		// Send push notification
		sendPushNotification(msg.UserID, msg.Body)
	}
}

func sendPushNotification(userID, body string) {
	var token string
	err := db.QueryRow("SELECT token FROM device_tokens WHERE user_id = $1", userID).Scan(&token)
	if err != nil {
		log.Printf("Failed to retrieve token for user %s: %v", userID, err)
		return
	}

	// Here you would integrate with FCM to send the notification using the token
	log.Printf("Sending push notification to user %s with token %s: %s", userID, token, body)
}

func main() {
	initDB()

	r := gin.Default()

	r.POST("/register", registerUser)
	r.POST("/register-token", registerDeviceToken)
	r.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})

	go handleMessages()

	port := "8080"
	fmt.Printf("Server running on http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
