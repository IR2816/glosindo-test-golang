package controllers

import (
	"fmt"
	"net/http"
	"time"

	"glosindo-backend-go/database"
	"glosindo-backend-go/models"
	"glosindo-backend-go/utils"

	"github.com/gin-gonic/gin"
)

type CheckInRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Address   string  `json:"address" binding:"required"`
}

type CheckOutRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Address   string  `json:"address" binding:"required"`
}

func CheckIn(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	var req CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi GPS
	if !utils.IsInOfficeArea(req.Latitude, req.Longitude) {
		distance := utils.GetDistanceFromOffice(req.Latitude, req.Longitude)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Anda berada di luar area kantor (%.0fm dari kantor)", distance),
		})
		return
	}

	// Check apakah sudah check-in hari ini
	today := time.Now().Format("2006-01-02")
	var existingPresensi models.Presensi
	err := database.DB.Where("user_id = ? AND DATE(date) = ?", user.ID, today).First(&existingPresensi).Error

	if err == nil && existingPresensi.CheckInTime != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan check-in hari ini"})
		return
	}

	now := time.Now()
	dateOnly, _ := time.Parse("2006-01-02", today)

	if err == nil {
		// Update existing record
		existingPresensi.CheckInTime = &now
		existingPresensi.CheckInLat = &req.Latitude
		existingPresensi.CheckInLng = &req.Longitude
		existingPresensi.CheckInAddress = req.Address
		existingPresensi.Status = "hadir"

		if err := database.DB.Save(&existingPresensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update check-in"})
			return
		}

		c.JSON(http.StatusOK, existingPresensi)
		return
	}

	// Create new record
	presensi := models.Presensi{
		UserID:         user.ID,
		Date:           dateOnly,
		CheckInTime:    &now,
		CheckInLat:     &req.Latitude,
		CheckInLng:     &req.Longitude,
		CheckInAddress: req.Address,
		Status:         "hadir",
	}

	if err := database.DB.Create(&presensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create check-in"})
		return
	}

	c.JSON(http.StatusOK, presensi)
}

func CheckOut(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	var req CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check apakah sudah check-in
	today := time.Now().Format("2006-01-02")
	var presensi models.Presensi
	err := database.DB.Where("user_id = ? AND DATE(date) = ?", user.ID, today).First(&presensi).Error

	if err != nil || presensi.CheckInTime == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum melakukan check-in hari ini"})
		return
	}

	if presensi.CheckOutTime != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda sudah melakukan check-out hari ini"})
		return
	}

	// Update check-out
	now := time.Now()
	presensi.CheckOutTime = &now
	presensi.CheckOutLat = &req.Latitude
	presensi.CheckOutLng = &req.Longitude
	presensi.CheckOutAddress = req.Address

	if err := database.DB.Save(&presensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update check-out"})
		return
	}

	c.JSON(http.StatusOK, presensi)
}

func GetTodayPresensi(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	today := time.Now().Format("2006-01-02")
	var presensi models.Presensi
	err := database.DB.Where("user_id = ? AND DATE(date) = ?", user.ID, today).First(&presensi).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Presensi hari ini belum ada"})
		return
	}

	c.JSON(http.StatusOK, presensi)
}

func GetPresensiHistory(c *gin.Context) {
	currentUser, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	user := currentUser.(models.User)

	limit := 30
	if l := c.Query("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}

	var presensiList []models.Presensi
	if err := database.DB.Where("user_id = ?", user.ID).
		Order("date DESC").
		Limit(limit).
		Find(&presensiList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch presensi history"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": presensiList,
	})
}

func GetPresensiStats(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	now := time.Now()
	month := now.Month()
	year := now.Year()

	if m := c.Query("month"); m != "" {
		fmt.Sscanf(m, "%d", &month)
	}
	if y := c.Query("year"); y != "" {
		fmt.Sscanf(y, "%d", &year)
	}

	var presensiList []models.Presensi
	database.DB.Where("user_id = ? AND EXTRACT(MONTH FROM date) = ? AND EXTRACT(YEAR FROM date) = ?",
		user.ID, int(month), year).Find(&presensiList)

	totalHadir := 0
	totalTerlambat := 0
	totalIzin := 0
	totalAlpha := 0

	for _, p := range presensiList {
		if p.Status == "hadir" {
			totalHadir++
			// Check terlambat (setelah jam 09:00)
			if p.CheckInTime != nil && p.CheckInTime.Hour() >= 9 {
				totalTerlambat++
			}
		} else if p.Status == "izin" || p.Status == "sakit" {
			totalIzin++
		} else if p.Status == "alpha" {
			totalAlpha++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_hadir":     totalHadir,
		"total_terlambat": totalTerlambat,
		"total_izin":      totalIzin,
		"total_alpha":     totalAlpha,
	})
}