// Login and signup logic
package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ArpitKRS/go-med-record/config"
	"github.com/ArpitKRS/go-med-record/models"
	"github.com/ArpitKRS/go-med-record/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var usersCollection *mongo.Collection

// Initialize the users collection
func init() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	usersCollection = db.Collection("users")
}

// Signup handles new user registration
func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Insert user document into the collection
	_, err = usersCollection.InsertOne(config.Ctx, user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User registered successfully")
}

// Login handles user authentication and token generation
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Aadhar   string `json:"aadhar"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var user models.User
	err := usersCollection.FindOne(config.Ctx, bson.M{"aadhar": credentials.Aadhar}).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.Hex(), user.Role)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
