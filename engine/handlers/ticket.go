package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"war-ticket-engine/models"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
)

var (
	// Location quotas (atomic)
	locationQuotas = map[string]*int64{
		"graha-dipta":   new(int64),
		"juanda":        new(int64),
		"gedung-antam":  new(int64),
		"setiabudi-one": new(int64),
	}

	locationNames = map[string]string{
		"graha-dipta":   "Butik Emas LM - Graha Dipta",
		"juanda":        "Butik Emas LM - Juanda",
		"gedung-antam":  "Butik Emas LM - Gedung Antam",
		"setiabudi-one": "Butik Emas LM - Setiabudi One",
	}

	ticketCounter int64 = 0
	ticketsMu     sync.RWMutex
)

func init() {
	// Initialize quotas
	*locationQuotas["graha-dipta"] = 30
	*locationQuotas["juanda"] = 25
	*locationQuotas["gedung-antam"] = 40
	*locationQuotas["setiabudi-one"] = 20
}

func generateTicketNumber() string {
	num := atomic.AddInt64(&ticketCounter, 1)
	letters := []byte("ABCDEFGHJKLMNPQRSTUVWXYZ")
	prefix := string(letters[num%24]) + string(letters[(num/24)%24]) + string(letters[(num/576)%24])
	return fmt.Sprintf("%s%d-%03d", prefix, (num/1000)%10, num%1000)
}

func generateCode() string {
	bytes := make([]byte, 3)
	rand.Read(bytes)
	return fmt.Sprintf("%s%s", hex.EncodeToString(bytes)[:3], hex.EncodeToString(bytes)[3:])
}

type CreateTicketRequest struct {
	UserID     string `json:"user_id"`
	LocationID string `json:"location_id"`
	TimeSlot   string `json:"time_slot"`
}

func CreateTicketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
			return
		}

		quotaPtr, exists := locationQuotas[req.LocationID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Lokasi tidak valid"})
			return
		}

		// Atomic decrease
		remaining := atomic.AddInt64(quotaPtr, -1)
		if remaining < 0 {
			atomic.AddInt64(quotaPtr, 1) // Rollback
			c.JSON(http.StatusOK, gin.H{
				"status":  "failed",
				"message": "Kuota habis untuk lokasi ini",
			})
			return
		}

		ticket := models.Ticket{
			ID:           generateID(),
			UserID:       req.UserID,
			LocationID:   req.LocationID,
			LocationName: locationNames[req.LocationID],
			TicketNumber: generateTicketNumber(),
			Code:         strings.ToUpper(generateCode()),
			TimeSlot:     req.TimeSlot,
			CreatedAt:    time.Now(),
		}

		services.DB.SetTicket(ticket)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Tiket berhasil dibuat",
			"ticket":  ticket,
		})
	}
}

func GetLocationsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		locations := []gin.H{}
		for id, name := range locationNames {
			quota := atomic.LoadInt64(locationQuotas[id])
			locations = append(locations, gin.H{
				"id":    id,
				"name":  name,
				"quota": quota,
			})
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	}
}

func GetTicketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ticketID := c.Param("id")

		ticket, exists := services.DB.GetTicket(ticketID)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Tiket tidak ditemukan"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success", "ticket": ticket})
	}
}
