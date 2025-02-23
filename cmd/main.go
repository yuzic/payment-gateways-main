package main

import (
	"log"
	"net/http"
	"os"

	"payment-gateway/db"
	"payment-gateway/internal/api"
	"payment-gateway/internal/kafka"
)

func main() {
	// Initialize the database connection
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dbURL := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=disable"

	dbConnect, err := db.InitializeDB(dbURL)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successfully connected to the database.")

	// func init() {
	kafkaURL := os.Getenv("KAFKA_BROKER_URL")
	if kafkaURL == "" {
		kafkaURL = "kafka:9092"
	}

	kafkaPublisher := kafka.NewPublisher(kafkaURL)
	defer kafkaPublisher.Close()

	log.Println("Kafka writer initialized successfully.")
	// Set up the HTTP server and routes
	di := api.GetContainer(dbConnect, kafkaPublisher)
	router := api.SetupRouter(di)

	// Start the server on port 8080
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}

}
