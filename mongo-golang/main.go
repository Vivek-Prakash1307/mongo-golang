package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Vivek-Prakash1307/mongo-golang/controllers"
)

// Function to connect to MongoDB
func connectDB() (*mongo.Client, context.Context) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection was successful
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	} else {
		fmt.Println("Connected to MongoDB!")
	}

	return client, context.Background()
}

func main() {
	client, ctx := connectDB()
	defer client.Disconnect(ctx) // Close connection when program ends

	// Initialize router
	router := httprouter.New()

	// Pass MongoDB client to UserController
	uc := controllers.NewUserController(client)

	// Define API routes
	router.POST("/user", uc.CreateUser)       // Create user
	router.GET("/user/:id", uc.GetUser)       // Get user by ID
	router.DELETE("/user/:id", uc.DeleteUser) // Delete user by ID

	// Start server
	fmt.Println("Server running on port 9000...")
	log.Fatal(http.ListenAndServe(":9000", router))
}
