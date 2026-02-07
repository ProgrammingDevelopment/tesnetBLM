package handlers

import (
	"fmt"
	"net/http"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
)

type WarRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func WarHandler(redis *services.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req WarRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Atomic DECR
		remaining, err := redis.AtomicDecreaseQuota()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Redis error"})
			return
		}

		if remaining < 0 {
			// Rollback (optional, or just ignore since it's < 0)
			// For high perf, we might not INCR back immediately to avoid thrashing,
			// but effectively stock is empty.
			c.JSON(http.StatusOK, gin.H{
				"status":    "failed",
				"message":   "Quota habis",
				"remaining": 0,
			})
			return
		}

		// Success logic - Queue to DB (Simulated here)
		// go database.SaveTicket(...)

		c.JSON(http.StatusOK, gin.H{
			"status":        "success",
			"ticket_number": fmt.Sprintf("T-%d", 5000-remaining),
			"remaining":     remaining,
		})
	}
}

func StatusHandler(redis *services.RedisService) gin.HandlerFunc {
	return func(c *gin.Context) {
		count, _ := redis.GetQuota()
		c.JSON(http.StatusOK, gin.H{
			"quota_remaining": count,
		})
	}
}
