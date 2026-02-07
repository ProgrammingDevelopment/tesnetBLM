package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"net/mail"
	"regexp"
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

func validateRegisterInput(req RegisterRequest) error {
	if len(req.NIK) != 16 || !isAllDigits(req.NIK) {
		return errors.New("NIK harus 16 digit angka")
	}
	if !isValidPhone(req.Whatsapp) {
		return errors.New("No. WhatsApp tidak valid (contoh: 08xxxxxxxx)")
	}
	if !isValidEmail(req.Email) {
		return errors.New("Email tidak valid")
	}
	if !isStrongPassword(req.Password) {
		return errors.New("Password harus min 8 karakter dan mengandung huruf besar, huruf kecil, dan angka")
	}
	return nil
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return s != ""
}

func isValidPhone(p string) bool {
	// Basic Indonesia mobile pattern: starts with 08 and 10-15 digits
	if len(p) < 10 || len(p) > 15 {
		return false
	}
	if !strings.HasPrefix(p, "08") {
		return false
	}
	return isAllDigits(p)
}

func isValidEmail(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

func isStrongPassword(pw string) bool {
	if len(pw) < 8 {
		return false
	}
	var hasUpper, hasLower, hasDigit bool
	for _, r := range pw {
		switch {
		case r >= 'A' && r <= 'Z':
			hasUpper = true
		case r >= 'a' && r <= 'z':
			hasLower = true
		case r >= '0' && r <= '9':
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func sanitizeName(name string) string {
	// keep letters, spaces, dots, and hyphens
	re := regexp.MustCompile(`[^A-Za-z .-]+`)
	clean := re.ReplaceAllString(name, "")
	return strings.TrimSpace(clean)
}

func RegisterHandler(captcha *services.CaptchaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request"})
			return
		}

		if err := validateRegisterInput(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
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
			Nama:      strings.ToUpper(sanitizeName(req.Nama)),
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

		if req.Identifier == "" || req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Identifier dan password wajib diisi"})
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
