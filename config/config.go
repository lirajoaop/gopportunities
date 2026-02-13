package config

import (
	"fmt"

	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	logger *Logger
)

func Init() error {
	var err error

	// Initialize Database
	db, err = InitializeDatabase()

	if err != nil {
		return fmt.Errorf("error initializing database: %v", err)
	}

	return nil
}

func GetDatabase() *gorm.DB {
	return db
}

func GetLogger(prefix string) *Logger {
	// Initialize logger
	logger = NewLogger(prefix)
	return logger
}
