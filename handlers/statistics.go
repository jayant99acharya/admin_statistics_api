package handlers

import (
	"net/http"
	"time"

	"admin_statistics_api/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatisticsHandler struct {
	service   *services.StatisticsService
	validator *validator.Validate
}

type TimeRangeQuery struct {
	From string `form:"from" validate:"required" binding:"required"`
	To   string `form:"to" validate:"required" binding:"required"`
}

func NewStatisticsHandler() *StatisticsHandler {
	return &StatisticsHandler{
		service:   services.NewStatisticsService(),
		validator: validator.New(),
	}
}

func (h *StatisticsHandler) parseTimeRange(c *gin.Context) (time.Time, time.Time, error) {
	var query TimeRangeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Parse from date
	from, err := time.Parse("2006-01-02", query.From)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Parse to date
	to, err := time.Parse("2006-01-02", query.To)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Set to end of day for 'to' date
	to = to.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// Validate date range
	if from.After(to) {
		return time.Time{}, time.Time{}, gin.Error{
			Err:  nil,
			Type: gin.ErrorTypePublic,
			Meta: "from date cannot be after to date",
		}
	}

	// Validate dates are not in the future
	now := time.Now()
	if from.After(now) || to.After(now) {
		return time.Time{}, time.Time{}, gin.Error{
			Err:  nil,
			Type: gin.ErrorTypePublic,
			Meta: "dates cannot be in the future",
		}
	}

	return from, to, nil
}

// GetGrossGamingRevenue handles GET /gross_gaming_rev
func (h *StatisticsHandler) GetGrossGamingRevenue(c *gin.Context) {
	from, to, err := h.parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date parameters",
			"details": err.Error(),
		})
		return
	}

	results, err := h.service.GetGrossGamingRevenue(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate gross gaming revenue",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"from":                 from.Format("2006-01-02"),
			"to":                   to.Format("2006-01-02"),
			"gross_gaming_revenue": results,
		},
	})
}

// GetDailyWagerVolume handles GET /daily_wager_volume
func (h *StatisticsHandler) GetDailyWagerVolume(c *gin.Context) {
	from, to, err := h.parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date parameters",
			"details": err.Error(),
		})
		return
	}

	results, err := h.service.GetDailyWagerVolume(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate daily wager volume",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"from":                from.Format("2006-01-02"),
			"to":                  to.Format("2006-01-02"),
			"daily_wager_volume": results,
		},
	})
}

// GetUserWagerPercentile handles GET /user/:user_id/wager_percentile
func (h *StatisticsHandler) GetUserWagerPercentile(c *gin.Context) {
	// Parse user ID
	userIDStr := c.Param("user_id")
	userID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format",
			"details": "User ID must be a valid MongoDB ObjectID",
		})
		return
	}

	from, to, err := h.parseTimeRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid date parameters",
			"details": err.Error(),
		})
		return
	}

	result, err := h.service.GetUserWagerPercentile(userID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to calculate user wager percentile",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"from":             from.Format("2006-01-02"),
			"to":               to.Format("2006-01-02"),
			"user_percentile": result,
		},
	})
}

// HealthCheck handles GET /health
func (h *StatisticsHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": time.Now().UTC(),
		"service":   "Admin Statistics API",
	})
}