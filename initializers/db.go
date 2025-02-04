package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Declare the global DB variable

func ConnectDB() {
	// Load the Data Source Name (DSN) from the environment variable
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		log.Fatal("Environment variable DB_URL is not set")
	}

	// Connect to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Assign to the global variable

	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Connected to the database successfully")
}
