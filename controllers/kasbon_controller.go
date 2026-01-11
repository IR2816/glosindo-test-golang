package controllers

import (
	"net/http"

	"glosindo-backend-go/database"
	"glosindo-backend-go/models"

	"github.com/gin-gonic/gin"
)

type CreateKasbonRequest struct {
	Nominal float64 `json:"nominal" binding:"required,gt=0"`
	Reason  string  `json:"reason" binding:"required"`
}

// Get all kasbon
func GetKasbon(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	status := c.Query("status")

	query := database.DB.Where("user_id = ?", user.ID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var kasbonList []models.Kasbon
	query.Order("created_at DESC").Find(&kasbonList)

	c.JSON(http.StatusOK, kasbonList)
}

// Create kasbon
func CreateKasbon(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	var req CreateKasbonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	kasbon := models.Kasbon{
		UserID:  user.ID,
		Nominal: req.Nominal,
		Reason:  req.Reason,
		Status:  "pending",
	}

	if err := database.DB.Create(&kasbon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create kasbon"})
		return
	}

	c.JSON(http.StatusOK, kasbon)
}

// Get kasbon statistics
func GetKasbonStats(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	var stats struct {
		TotalPending  int64   `json:"total_pending"`
		TotalApproved int64   `json:"total_approved"`
		TotalRejected int64   `json:"total_rejected"`
		TotalNominal  float64 `json:"total_nominal"`
	}

	database.DB.Model(&models.Kasbon{}).Where("user_id = ? AND status = ?", user.ID, "pending").Count(&stats.TotalPending)
	database.DB.Model(&models.Kasbon{}).Where("user_id = ? AND status = ?", user.ID, "approved").Count(&stats.TotalApproved)
	database.DB.Model(&models.Kasbon{}).Where("user_id = ? AND status = ?", user.ID, "rejected").Count(&stats.TotalRejected)

	database.DB.Model(&models.Kasbon{}).
		Where("user_id = ? AND status = ?", user.ID, "approved").
		Select("COALESCE(SUM(nominal), 0)").
		Scan(&stats.TotalNominal)

	c.JSON(http.StatusOK, stats)
}