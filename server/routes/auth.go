package routes

import (
	"github.com/ArpitKRS/go-med-record/controllers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/api/login", controllers.Login).Methods("POST")

	documentRoutes := router.PathPrefix("/api").Subrouter()
	documentRoutes.HandleFunc("/upload", controllers.UploadDocument).Methods("POST")
	documentRoutes.HandleFunc("/user/documents", controllers.ViewOwnDocuments).Methods("GET")
	documentRoutes.HandleFunc("/doctor/medical-history", controllers.ViewPatientDocuments).Methods("GET")
	documentRoutes.Use(AuthMiddleware)

	return router
}
