package handlers

import (
	"net/http"
	"war-ticket-engine/services"

	"github.com/gin-gonic/gin"
)

func MathCaptchaHandler(captcha *services.CaptchaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		challenge, err := captcha.NewMathCaptcha()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Captcha error"})
			return
		}
		c.JSON(http.StatusOK, challenge)
	}
}

func ImageCaptchaHandler(captcha *services.CaptchaService) gin.HandlerFunc {
	return func(c *gin.Context) {
		challenge, err := captcha.NewImageCaptcha()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Captcha error"})
			return
		}
		c.JSON(http.StatusOK, challenge)
	}
}
