package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	// Get database configuration from environment variables or use defaults
	dbUser := getEnv("DB_USER", "root")
	dbPass := os.Getenv("DB_PASSWORD") // Get without fallback to detect if really empty
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306") // MySQL default port
	dbName := getEnv("DB_NAME", "simnikah")

	// Log configuration (WITHOUT PASSWORD for security!)
	log.Printf("=== DATABASE CONFIGURATION ===")
	log.Printf("DB_HOST: %s", dbHost)
	log.Printf("DB_PORT: %s", dbPort)
	log.Printf("DB_USER: %s", dbUser)
	log.Printf("DB_NAME: %s", dbName)
	log.Printf("DB_PASSWORD: %s", maskPassword(dbPass))
	log.Printf("==============================")

	// Validate required environment variables
	if dbPass == "" {
		// Allow empty password for local development (localhost/127.0.0.1)
		if dbHost != "127.0.0.1" && dbHost != "localhost" {
			log.Printf("ERROR: DB_PASSWORD is EMPTY!")
			log.Printf("Railway Troubleshooting:")
			log.Printf("1. Check if MySQL service is added to your Railway project")
			log.Printf("2. In your application service, go to Variables tab")
			log.Printf("3. Make sure DB_PASSWORD variable is set (not empty)")
			log.Printf("4. Try setting it manually by copying MYSQL_ROOT_PASSWORD from MySQL service")
			return nil, fmt.Errorf("DB_PASSWORD environment variable is required for non-localhost connections (current host: %s)", dbHost)
		}
		log.Println("⚠️  Warning: Using empty password for localhost MySQL")
	}
	if dbHost == "127.0.0.1" || dbHost == "localhost" {
		log.Println("⚠️  Warning: Using localhost database. Make sure to set DB_HOST for production.")
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

	// Create DSN string for MySQL
	// Format: user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// Use UTC timezone for compatibility (both local and Railway)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// Log connection details (without password)
	log.Printf("Connecting to MySQL: host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL database: %v", err)
	}

	log.Println("Connected to MySQL database successfully")

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	// Set connection pool settings optimized for Railway
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection pool configured successfully")

	return db, nil
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Helper function to mask password for logging (security!)
func maskPassword(password string) string {
	if password == "" {
		return "[EMPTY]"
	}
	if len(password) <= 4 {
		return "****"
	}
	// Show first 2 and last 2 characters, mask the rest
	return password[:2] + strings.Repeat("*", len(password)-4) + password[len(password)-2:]
}
