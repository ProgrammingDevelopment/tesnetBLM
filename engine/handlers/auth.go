package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
	"war-ticket-engine/models"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	NIK                string `json:"nik"`
	Nama               string `json:"nama"`
	Whatsapp           string `json:"whatsapp"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	CaptchaMathToken   string `json:"captcha_math_token"`
	CaptchaMathAnswer  string `json:"captcha_math_answer"`
	CaptchaImageToken  string `json:"captcha_image_token"`
	CaptchaImageAnswer string `json:"captcha_image_answer"`
}

type LoginRequest struct {
	Identifier         string `json:"identifier"` // email or whatsapp
	Password           string `json:"password"`
	CaptchaMathToken   string `json:"captcha_math_token"`
	CaptchaMathAnswer  string `json:"captcha_math_answer"`
	CaptchaImageToken  string `json:"captcha_image_token"`
	CaptchaImageAnswer string `json:"captcha_image_answer"`
}

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func RegisterHandler(captcha *services.CaptchaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
			return
		}

		if captcha != nil {
			if err := captcha.ValidateMath(req.CaptchaMathToken, req.CaptchaMathAnswer); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Captcha tidak valid"})
				return
			}
			if err := captcha.ValidateImage(req.CaptchaImageToken, req.CaptchaImageAnswer); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Captcha tidak valid"})
				return
			}
		}

		// Validate NIK
		if len(req.NIK) != 16 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "NIK harus 16 digit"})
			return
		}

		// Check existing user
		if _, exists := services.DB.GetUserByNIK(req.NIK); exists {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "NIK sudah terdaftar"})
			return
		}
		if _, exists := services.DB.GetUserByEmailOrPhone(req.Email); exists {
			c.JSON(http.StatusConflict, gin.H{"status": "error", "message": "Email sudah terdaftar"})
			return
		}

		// Hash password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

		user := models.User{
			ID:        generateID(),
			NIK:       req.NIK,
			Nama:      strings.ToUpper(req.Nama),
			Whatsapp:  req.Whatsapp,
			Email:     req.Email,
			Password:  string(hashedPassword),
			CreatedAt: time.Now(),
		}

		services.DB.SetUser(user)

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Registrasi berhasil",
			"user": gin.H{
				"id":       user.ID,
				"nik":      user.NIK,
				"nama":     user.Nama,
				"whatsapp": user.Whatsapp,
				"email":    user.Email,
			},
		})
	}
}

func LoginHandler(captcha *services.CaptchaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
			return
		}

		if captcha != nil {
			if err := captcha.ValidateMath(req.CaptchaMathToken, req.CaptchaMathAnswer); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Captcha tidak valid"})
				return
			}
			if err := captcha.ValidateImage(req.CaptchaImageToken, req.CaptchaImageAnswer); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Captcha tidak valid"})
				return
			}
		}

		user, exists := services.DB.GetUserByEmailOrPhone(req.Identifier)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "User tidak ditemukan"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Password salah"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Login berhasil",
			"user": gin.H{
				"id":       user.ID,
				"nik":      user.NIK,
				"nama":     user.Nama,
				"whatsapp": user.Whatsapp,
				"email":    user.Email,
			},
		})
	}
}
