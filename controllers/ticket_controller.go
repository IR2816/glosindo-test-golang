package controllers

import (
	"net/http"
	"time"

	"glosindo-backend-go/database"
	"glosindo-backend-go/models"

	"github.com/gin-gonic/gin"
)

type CreateTicketRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Priority    string `json:"priority"`
}

type UpdateTicketStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Notes  string `json:"notes"`
}

// Get all tickets with filters
func GetTickets(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	// Query params
	status := c.Query("status")
	category := c.Query("category")
	search := c.Query("search")

	query := database.DB.Where("user_id = ?", user.ID)

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var tickets []models.Ticket
	query.Preload("Progress").Order("created_at DESC").Find(&tickets)

	c.JSON(http.StatusOK, tickets)
}

// Get single ticket detail
func GetTicketDetail(c *gin.Context) {
	ticketID := c.Param("id")

	var ticket models.Ticket
	if err := database.DB.Preload("Progress").Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// Create new ticket
func CreateTicket(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	var req CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	ticket := models.Ticket{
		UserID:      user.ID,
		Title:       req.Title,
		Description: req.Description,
		Category:    req.Category,
		Priority:    priority,
		Status:      "open",
	}

	if err := database.DB.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create ticket"})
		return
	}

	// Create initial progress
	progress := models.TicketProgress{
		TicketID:  ticket.ID,
		Status:    "open",
		Notes:     "Ticket dibuat",
		CreatedBy: user.ID,
	}
	database.DB.Create(&progress)

	// Load progress
	database.DB.Preload("Progress").First(&ticket, ticket.ID)

	c.JSON(http.StatusOK, ticket)
}

// Update ticket status
func UpdateTicketStatus(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)
	ticketID := c.Param("id")

	var req UpdateTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ticket models.Ticket
	if err := database.DB.Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Update ticket status
	ticket.Status = req.Status
	if req.Status == "closed" {
		now := time.Now()
		ticket.ClosedAt = &now
	}

	database.DB.Save(&ticket)

	// Create progress entry
	progress := models.TicketProgress{
		TicketID:  ticket.ID,
		Status:    req.Status,
		Notes:     req.Notes,
		CreatedBy: user.ID,
	}
	database.DB.Create(&progress)

	// Load progress
	database.DB.Preload("Progress").First(&ticket, ticket.ID)

	c.JSON(http.StatusOK, ticket)
}

// Get ticket statistics
func GetTicketStats(c *gin.Context) {
	currentUser, _ := c.Get("user")
	user := currentUser.(models.User)

	var stats struct {
		TotalOpen       int64 `json:"total_open"`
		TotalOnProgress int64 `json:"total_on_progress"`
		TotalClosed     int64 `json:"total_closed"`
		TotalTickets    int64 `json:"total_tickets"`
	}

	database.DB.Model(&models.Ticket{}).Where("user_id = ? AND status = ?", user.ID, "open").Count(&stats.TotalOpen)
	database.DB.Model(&models.Ticket{}).Where("user_id = ? AND status = ?", user.ID, "on_progress").Count(&stats.TotalOnProgress)
	database.DB.Model(&models.Ticket{}).Where("user_id = ? AND status = ?", user.ID, "closed").Count(&stats.TotalClosed)
	database.DB.Model(&models.Ticket{}).Where("user_id = ?", user.ID).Count(&stats.TotalTickets)

	c.JSON(http.StatusOK, stats)
}

// Get categories
func GetTicketCategories(c *gin.Context) {
	categories := []string{
		"IT Support",
		"Bug Report",
		"Feature Request",
		"Access Request",
		"Hardware Issue",
		"Software Issue",
		"Network Issue",
		"Other",
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}