package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL       string
	JWTSecret         string
	JWTExpireHours    int
	OfficeLatitude    float64
	OfficeLongitude   float64
	MaxDistanceMeters float64
	Port              string
	CORSOrigins       []string
}

var AppConfig *Config

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbURL := getEnv("DATABASE_URL", "postgresql://postgres:luVrWlXxolHaPBGmPqUdDSuXyaijkAMW@postgres.railway.internal:5432/railway")
	log.Printf("📋 Database URL configured (masked for security)")

	AppConfig = &Config{
		DatabaseURL:       dbURL,
		JWTSecret:         getEnv("JWT_SECRET", "glosindo_secret_key"),
		JWTExpireHours:    getEnvAsInt("JWT_EXPIRE_HOURS", 168),
		OfficeLatitude:    getEnvAsFloat("OFFICE_LATITUDE", -6.5947),
		OfficeLongitude:   getEnvAsFloat("OFFICE_LONGITUDE", 106.7890),
		MaxDistanceMeters: getEnvAsFloat("MAX_DISTANCE_METERS", 100),
		Port:              getEnv("PORT", "8000"),
		CORSOrigins:       []string{"http://localhost:52302", "https://your-frontend-domain.com"},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultValue
}
