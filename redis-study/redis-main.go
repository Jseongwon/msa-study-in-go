package redis_study

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
)

// Global Redis client
var redisClient *redis.Client
var ctx = context.Background()

// WebSocket connections
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Redis channels
const chatChannel = "chatroom"

func initRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Check Redis connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis")
}

// WebSocket handler
func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			delete(clients, ws)
			break
		}

		// Publish message to Redis
		err = redisClient.Publish(ctx, chatChannel, string(msg)).Err()
		if err != nil {
			log.Println("Publish error:", err)
		}
	}
}

// Redis subscriber for real-time chat
func redisSubscriber() {
	pubsub := redisClient.Subscribe(ctx, chatChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("Received from Redis: %s\n", msg.Payload)
		broadcast <- msg.Payload
	}
}

// Broadcast messages to WebSocket clients
func broadcaster() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Broadcast error:", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func main() {
	initRedis()

	// Start Redis subscriber in a separate goroutine
	go redisSubscriber()

	// Start WebSocket broadcaster in a separate goroutine
	go broadcaster()

	// Setup Gin server
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})

	fmt.Println("Chat server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
