package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Vivek-Prakash1307/mongo-golang/models"
)

type UserController struct {
	client *mongo.Client
}

// NewUserController creates a new UserController
func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

// Create User
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	collection := uc.client.Database("mongo-golang").Collection("users")

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate a new ObjectID
	user.Id = primitive.NewObjectID()

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Error inserting user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// Get User by ID
func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	collection := uc.client.Database("mongo-golang").Collection("users")
	id := p.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ObjectID format", http.StatusBadRequest)
		return
	}

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found in database", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// Delete User
func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	collection := uc.client.Database("mongo-golang").Collection("users")
	id := p.ByName("id")

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	_, err = collection.DeleteOne(context.Background(), bson.M{"_id": objID})
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s deleted successfully", id)
}
