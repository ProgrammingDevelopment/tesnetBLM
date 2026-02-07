package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"war-ticket-engine/models"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
)

var (
	// Location quotas (atomic) and metadata
	locationData = map[string]struct {
		Name   string
		Region string
		Quota  *int64
	}{
		"graha-dipta":   {Name: "Butik Emas LM - Graha Dipta", Region: "jabodetabek", Quota: new(int64)},
		"juanda":        {Name: "Butik Emas LM - Juanda", Region: "jabodetabek", Quota: new(int64)},
		"gedung-antam":  {Name: "Butik Emas LM - Gedung Antam", Region: "jabodetabek", Quota: new(int64)},
		"setiabudi-one": {Name: "Butik Emas LM - Setiabudi One", Region: "jabodetabek", Quota: new(int64)},
	}

	ticketCounter int64 = 0
	ticketsMu     sync.RWMutex
)

func init() {
	// Initialize quotas
	*locationData["graha-dipta"].Quota = 30
	*locationData["juanda"].Quota = 25
	*locationData["gedung-antam"].Quota = 40
	*locationData["setiabudi-one"].Quota = 20
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
	UserID     string  `json:"user_id"`
	LocationID string  `json:"location_id"`
	TimeSlot   string  `json:"time_slot"`
	SizeGram   float64 `json:"size_gram,omitempty"`
}

func CreateTicketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
			return
		}

		location, exists := locationData[req.LocationID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Lokasi tidak valid"})
			return
		}

		// Atomic decrease
		remaining := atomic.AddInt64(location.Quota, -1)
		if remaining < 0 {
			atomic.AddInt64(location.Quota, 1) // Rollback
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
			LocationName: location.Name,
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
		for id, loc := range locationData {
			quota := atomic.LoadInt64(loc.Quota)
			locations = append(locations, gin.H{
				"id":     id,
				"name":   loc.Name,
				"region": loc.Region,
				"quota":  quota,
			})
		}
		c.JSON(http.StatusOK, gin.H{"locations": locations})
	}
}

type PreOpenTicketRequest struct {
	UserID     string  `json:"user_id"`
	LocationID string  `json:"location_id"`
	TimeSlot   string  `json:"time_slot"`
	SizeGram   float64 `json:"size_gram"`
}

func PreOpenTicketHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !inPreOpenWindow() {
			c.JSON(http.StatusTooEarly, gin.H{"status": "error", "message": "Pre-open belum dibuka"})
			return
		}

		var req PreOpenTicketRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
			return
		}

		minSize := readMinPreOpenSize()
		if req.SizeGram < minSize {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": fmt.Sprintf("Minimal ukuran %.1f gram untuk pre-open", minSize)})
			return
		}

		loc, exists := locationData[req.LocationID]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Lokasi tidak valid"})
			return
		}
		if strings.ToLower(loc.Region) != "jabodetabek" {
			c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Lokasi hanya untuk area Jabodetabek selama pre-open"})
			return
		}

		remaining := atomic.AddInt64(loc.Quota, -1)
		if remaining < 0 {
			atomic.AddInt64(loc.Quota, 1)
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
			LocationName: loc.Name,
			TicketNumber: generateTicketNumber(),
			Code:         strings.ToUpper(generateCode()),
			TimeSlot:     req.TimeSlot,
			CreatedAt:    time.Now(),
		}

		services.DB.SetTicket(ticket)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Tiket pre-open berhasil dibuat",
			"ticket":  ticket,
		})
	}
}

func inPreOpenWindow() bool {
	startCfg := strings.TrimSpace(os.Getenv("GATE_START_TIME"))
	if startCfg == "" {
		startCfg = "07:00"
	}
	offsetMinutes := 10
	if raw := strings.TrimSpace(os.Getenv("PREOPEN_OFFSET_MINUTES")); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			offsetMinutes = v
		}
	}

	now := time.Now()
	todayStart, err := time.ParseInLocation("15:04", startCfg, now.Location())
	if err != nil {
		return false
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), todayStart.Hour(), todayStart.Minute(), 0, 0, now.Location())
	windowStart := start.Add(-time.Duration(offsetMinutes) * time.Minute)
	return now.After(windowStart) && now.Before(start)
}

func readMinPreOpenSize() float64 {
	raw := strings.TrimSpace(os.Getenv("MIN_PREOPEN_SIZE_GRAM"))
	if raw == "" {
		return 5.0
	}
	val, err := strconv.ParseFloat(raw, 64)
	if err != nil || val <= 0 {
		return 5.0
	}
	return val
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
