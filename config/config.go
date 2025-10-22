package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	// Get database configuration from environment variables or use defaults
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "")
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432") // Default PostgreSQL port, LeapCell uses 6438
	dbName := getEnv("DB_NAME", "simnikah")

	// Validate required environment variables
	if dbPass == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}
	if dbHost == "127.0.0.1" || dbHost == "localhost" {
		log.Println("Warning: Using localhost database. Make sure to set DB_HOST for production.")
	}

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Create DSN string for PostgreSQL
	// Try with SSL first, fallback to disable if needed
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Jakarta",
		dbHost, dbUser, dbPass, dbName, dbPort)

	// Log connection details (without password)
	log.Printf("Connecting to PostgreSQL: host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		// Try without SSL if SSL connection fails
		log.Printf("SSL connection failed, trying without SSL: %v", err)
		dsnNoSSL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			dbHost, dbUser, dbPass, dbName, dbPort)

		db, err = gorm.Open(postgres.Open(dsnNoSSL), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database with both SSL and non-SSL: %v", err)
		}
		log.Println("Connected to database without SSL")
	} else {
		log.Println("Connected to database with SSL")
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings optimized for LeapCell PostgreSQL
	sqlDB.SetMaxIdleConns(5)  // Reduced for 0.2 vCPU, 500MB RAM
	sqlDB.SetMaxOpenConns(20) // Reduced for limited resources
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
