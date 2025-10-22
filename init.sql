-- SimNikah Database Initialization Script
-- This script will be executed when MySQL container starts

-- Create database if not exists
CREATE DATABASE IF NOT EXISTS simnikah CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE simnikah;

-- Create initial admin user (optional)
-- This will be handled by the Go application's AutoMigrate

-- Set timezone
SET time_zone = '+07:00';

-- Create indexes for better performance (will be created by GORM AutoMigrate)
-- But we can add some custom indexes here if needed

-- Grant permissions to the application user
GRANT ALL PRIVILEGES ON simnikah.* TO 'simnikah_user'@'%';
FLUSH PRIVILEGES;
