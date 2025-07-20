package main

import (
	"log"
	"net/http"
	"os"
	_ "sub_service/docs"
	database "sub_service/internal/db"
	"sub_service/internal/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Subscription Service API
// @version 1.0
// @description API для управления подписками
// @host localhost:8080 (default)
// @BasePath /
// @schemes http
func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//setup port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	//db
	database.ConnectToPostgres()

	r := mux.NewRouter()
	// Routes
	r.HandleFunc("/subscriptions/CreateSub", routes.CreateSub).Methods("POST")
	r.HandleFunc("/subscriptions/UpdateSub/{id}", routes.SubUpdate).Methods("PUT")
	r.HandleFunc("/subscriptions/GetSub/{id}", routes.GetSub).Methods("GET")
	r.HandleFunc("/subscriptions/DeleteSub/{id}", routes.DeleteSub).Methods("DELETE")
	r.HandleFunc("/subscriptions/total-cost/{user_id}/{subscription_service}/{start_date_from}/{start_date_to}", routes.SubSum).Methods("GET")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Printf("Server started on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
