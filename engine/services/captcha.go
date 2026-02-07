package services

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"
)

type CaptchaService struct {
	secret    []byte
	ttl       time.Duration
	imageSets []imageChallenge
}

type MathCaptcha struct {
	A         int    `json:"a"`
	B         int    `json:"b"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type ImageCaptcha struct {
	Prompt    string   `json:"prompt"`
	Options   []string `json:"options"`
	Token     string   `json:"token"`
	ExpiresAt int64    `json:"expires_at"`
}

type imageChallenge struct {
	Prompt  string
	Correct string
	Options []string
}

func NewCaptchaService() *CaptchaService {
	secret := strings.TrimSpace(os.Getenv("CAPTCHA_SECRET"))
	if secret == "" {
		secret = "dev-secret"
	}

	ttlSeconds := int64(300)
	if raw := strings.TrimSpace(os.Getenv("CAPTCHA_TTL_SECONDS")); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err == nil && parsed > 0 {
			ttlSeconds = parsed
		}
	}

	return &CaptchaService{
		secret: []byte(secret),
		ttl:    time.Duration(ttlSeconds) * time.Second,
		imageSets: []imageChallenge{
			{
				Prompt:  "Pilih gambar emas batangan",
				Correct: "goldbar",
				Options: []string{"goldbar", "coin", "ring", "wallet"},
			},
			{
				Prompt:  "Pilih gambar koin",
				Correct: "coin",
				Options: []string{"goldbar", "coin", "ring", "wallet"},
			},
			{
				Prompt:  "Pilih gambar cincin",
				Correct: "ring",
				Options: []string{"goldbar", "coin", "ring", "wallet"},
			},
		},
	}
}

func (c *CaptchaService) NewMathCaptcha() (MathCaptcha, error) {
	a, err := randomInt(1, 9)
	if err != nil {
		return MathCaptcha{}, err
	}
	b, err := randomInt(1, 9)
	if err != nil {
		return MathCaptcha{}, err
	}

	exp := time.Now().Add(c.ttl).Unix()
	token := c.signToken(exp, "math", fmt.Sprintf("%d,%d", a, b))

	return MathCaptcha{
		A:         a,
		B:         b,
		Token:     token,
		ExpiresAt: exp,
	}, nil
}

func (c *CaptchaService) NewImageCaptcha() (ImageCaptcha, error) {
	if len(c.imageSets) == 0 {
		return ImageCaptcha{}, errors.New("no image captcha configured")
	}

	index, err := randomInt(0, len(c.imageSets)-1)
	if err != nil {
		return ImageCaptcha{}, err
	}
	challenge := c.imageSets[index]

	options := append([]string{}, challenge.Options...)
	shuffle(options)

	exp := time.Now().Add(c.ttl).Unix()
	token := c.signToken(exp, "image", challenge.Correct)

	return ImageCaptcha{
		Prompt:    challenge.Prompt,
		Options:   options,
		Token:     token,
		ExpiresAt: exp,
	}, nil
}

func (c *CaptchaService) ValidateMath(token, answer string) error {
	exp, kind, data, err := c.parseToken(token)
	if err != nil {
		return err
	}
	if kind != "math" {
		return errors.New("invalid captcha type")
	}
	if time.Now().Unix() > exp {
		return errors.New("captcha expired")
	}

	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return errors.New("invalid captcha data")
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		return errors.New("invalid captcha data")
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("invalid captcha data")
	}

	expected := a + b
	answer = strings.TrimSpace(answer)
	if answer == "" {
		return errors.New("captcha answer required")
	}
	val, err := strconv.Atoi(answer)
	if err != nil || val != expected {
		return errors.New("captcha answer incorrect")
	}
	return nil
}

func (c *CaptchaService) ValidateImage(token, answer string) error {
	exp, kind, data, err := c.parseToken(token)
	if err != nil {
		return err
	}
	if kind != "image" {
		return errors.New("invalid captcha type")
	}
	if time.Now().Unix() > exp {
		return errors.New("captcha expired")
	}

	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer == "" {
		return errors.New("captcha answer required")
	}
	if answer != strings.ToLower(data) {
		return errors.New("captcha answer incorrect")
	}
	return nil
}

func (c *CaptchaService) signToken(exp int64, kind, data string) string {
	payload := fmt.Sprintf("%d|%s|%s", exp, kind, data)
	mac := hmac.New(sha256.New, c.secret)
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	token := payload + "|" + sig
	return base64.RawURLEncoding.EncodeToString([]byte(token))
}

func (c *CaptchaService) parseToken(token string) (int64, string, string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(token))
	if err != nil {
		return 0, "", "", errors.New("invalid captcha token")
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 4 {
		return 0, "", "", errors.New("invalid captcha token")
	}
	exp, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", "", errors.New("invalid captcha token")
	}
	kind := parts[1]
	data := parts[2]
	sig := parts[3]

	payload := fmt.Sprintf("%d|%s|%s", exp, kind, data)
	mac := hmac.New(sha256.New, c.secret)
	mac.Write([]byte(payload))
	expected := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return 0, "", "", errors.New("invalid captcha token")
	}
	return exp, kind, data, nil
}

func randomInt(min, max int) (int, error) {
	if max < min {
		return 0, errors.New("invalid range")
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + min, nil
}

func shuffle(values []string) {
	for i := len(values) - 1; i > 0; i-- {
		j, err := randomInt(0, i)
		if err != nil {
			continue
		}
		values[i], values[j] = values[j], values[i]
	}
}
