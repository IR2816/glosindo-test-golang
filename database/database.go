package database

import (
	"fmt"
	"log"
	"strings"

	"glosindo-backend-go/config"
	"glosindo-backend-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	var dialector gorm.Dialector
	if strings.HasPrefix(config.AppConfig.DatabaseURL, "sqlite:") {
		dialector = sqlite.Open(strings.TrimPrefix(config.AppConfig.DatabaseURL, "sqlite:"))
	} else {
		dialector = postgres.Open(config.AppConfig.DatabaseURL)
	}

	DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("✅ Database connected")

	// Auto Migrate
	err = DB.AutoMigrate(
		&models.User{},
		&models.Presensi{},
		&models.LoginHistory{},
		&models.Ticket{},
		&models.TicketProgress{},
		&models.Kasbon{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("✅ Database migrated")
}
