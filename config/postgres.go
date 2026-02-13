package config

import (
	"fmt"

	"github.com/lirajoaop/gopportunities/schemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDatabase() (*gorm.DB, error) {
	logger := GetLogger("database")

	dsn := GetDatabaseURL()
	if dsn == "" {
		logger.Error("DATABASE_URL environment variable is not set")
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Errorf("database connection error: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&schemas.Opening{})
	if err != nil {
		logger.Errorf("database automigration error: %v", err)
		return nil, err
	}

	logger.Info("database connected and migrated successfully")
	return db, nil
}
