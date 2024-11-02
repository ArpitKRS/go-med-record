// To manage uploading of docs on the api
package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ArpitKRS/go-med-record/config"
	"github.com/ArpitKRS/go-med-record/models"
	"github.com/ArpitKRS/go-med-record/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var docCollection *mongo.Collection

// Initialize the documents collection
func init() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	docCollection = db.Collection("documents")
}

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	aadhar := r.FormValue("aadhar")
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File upload err", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	if ext != ".png" && ext != ".pdf" {
		http.Error(w, "Invalid file type (not png or pdf)", http.StatusBadRequest)
		return
	}

	// Save file to disk
	filePath := "uploads/" + handler.Filename
	os.MkdirAll("uploads", os.ModePerm) // Ensure uploads dir exists
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	hashedAadhar := utils.HashAadhar(aadhar)
	document := models.Document{
		UserID:  userID,
		Aadhar:  hashedAadhar,
		FileURL: filePath,
	}

	_, err = docCollection.InsertOne(context.TODO(), document)
	if err != nil {
		http.Error(w, "Error saving document", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Document successfully uploaded!")
}

func ViewOwnDocuments(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	cursor, err := docCollection.Find(context.TODO(), bson.M{"userID": userID})
	if err != nil {
		http.Error(w, "Error retrieving documents", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var documents []models.Document
	cursor.All(context.TODO(), &documents)

	json.NewEncoder(w).Encode(documents)
}

func ViewPatientDocuments(w http.ResponseWriter, r *http.Request) {
	userRole := r.Context().Value("role").(string)

	if userRole != "doctor" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	aadhar := r.URL.Query().Get("aadhar")
	hashedAadhar := utils.HashAadhar(aadhar)

	cursor, err := docCollection.Find(context.TODO(), bson.M{"aadhar": hashedAadhar})
	if err != nil {
		http.Error(w, "Error retrieving patient documents", http.StatusInternalServerError)
		return
	}

	defer cursor.Close(context.TODO())

	var documents []models.Document
	cursor.All(context.TODO(), &documents)

	json.NewEncoder(w).Encode(documents)
}
