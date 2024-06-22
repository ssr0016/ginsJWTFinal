package config

import (
	"fmt"
	"gins/helper"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectionDB(config *Config) *gorm.DB {
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DBHost, config.DBPort, config.DBUsername, config.DBPassword, config.DBName)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})
	helper.ErrorPanic(err)

	sqlDB, err := db.DB()
	helper.ErrorPanic(err)

	sqlDB.SetMaxOpenConns(config.DB.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.DB.MaxIdleConns)

	// Parse DB_MAX_IDLE_TIME
	durationStr := config.DB.MaxIdleTime
	if durationStr == "" {
		durationStr = "15m" // Default value
	}
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Log the error and set a default duration if the configuration value is invalid.
		fmt.Println("Invalid DB_MAX_IDLE_TIME, setting to default 15m:", err)
		duration = 15 * time.Minute
	}
	sqlDB.SetConnMaxIdleTime(duration)

	fmt.Println("Connected Successfully to the Database")
	return db
}
