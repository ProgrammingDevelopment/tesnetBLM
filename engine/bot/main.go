package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type statusResponse struct {
	QuotaRemaining int64 `json:"quota_remaining"`
}

type chatResponse struct {
	Answer string `json:"answer"`
}

func fetchBackendStatus(baseURL string) (bool, int64) {
	client := &http.Client{Timeout: 3 * time.Second}
	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(baseURL, "/")+"/api/status", nil)
	if err != nil {
		return false, 0
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, 0
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, 0
	}

	var payload statusResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return false, 0
	}

	return true, payload.QuotaRemaining
}

func fetchChatAnswer(baseURL, message string) (string, error) {
	client := &http.Client{Timeout: 20 * time.Second}
	payload := strings.NewReader(`{"message":"` + escapeJSON(message) + `"}`)
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(baseURL, "/")+"/api/chat", payload)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("chat request failed")
	}

	var out chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.Answer), nil
}

func escapeJSON(input string) string {
	replacer := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`)
	return replacer.Replace(input)
}

func handleUpdates(bot *tgbotapi.BotAPI, baseURL string, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "Welcome to War Tiket Bot. Use /status to check backend availability or ask a question."
			case "help":
				msg.Text = "Commands: /status, /help. You can also ask a question directly."
			case "status":
				online, quota := fetchBackendStatus(baseURL)
				if online {
					msg.Text = "Backend: ONLINE. Remaining quota: " + strconv.FormatInt(quota, 10)
				} else {
					msg.Text = "Backend: OFFLINE."
				}
			default:
				msg.Text = "Unknown command. Use /help."
			}
		} else if strings.TrimSpace(update.Message.Text) != "" {
			answer, err := fetchChatAnswer(baseURL, update.Message.Text)
			if err != nil || answer == "" {
				msg.Text = "Sorry, I could not reach the assistant service."
			} else {
				msg.Text = answer
			}
		} else {
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("send error: %v", err)
		}
	}
}

func startPolling(bot *tgbotapi.BotAPI, baseURL string) {
	_, _ = bot.Request(tgbotapi.DeleteWebhookConfig{DropPendingUpdates: true})
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	handleUpdates(bot, baseURL, updates)
}

func startWebhook(bot *tgbotapi.BotAPI, baseURL, webhookURL string) {
	parsed, err := url.Parse(webhookURL)
	if err != nil {
		log.Fatalf("Invalid TELEGRAM_WEBHOOK_URL: %v", err)
	}

	path := strings.TrimSpace(parsed.Path)
	if path == "" {
		path = strings.TrimSpace(os.Getenv("BOT_WEBHOOK_PATH"))
	}
	if path == "" {
		path = "/telegram/webhook"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if parsed.Path == "" {
		webhookURL = strings.TrimRight(webhookURL, "/") + path
	}

	wh, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		log.Fatalf("Webhook setup failed: %v", err)
	}
	if _, err := bot.Request(wh); err != nil {
		log.Fatalf("Webhook setup failed: %v", err)
	}

	listenAddr := strings.TrimSpace(os.Getenv("BOT_LISTEN_ADDR"))
	if listenAddr == "" {
		listenAddr = ":8081"
	}

	updates := bot.ListenForWebhook(path)
	go func() {
		log.Printf("Webhook listening on %s%s", listenAddr, path)
		if err := http.ListenAndServe(listenAddr, nil); err != nil {
			log.Fatalf("Webhook server error: %v", err)
		}
	}()

	handleUpdates(bot, baseURL, updates)
}

func main() {
	token := strings.TrimSpace(os.Getenv("TELEGRAM_APITOKEN"))
	if token == "" {
		log.Fatal("TELEGRAM_APITOKEN is required")
	}

	baseURL := strings.TrimSpace(os.Getenv("BACKEND_BASE_URL"))
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	webhookURL := strings.TrimSpace(os.Getenv("TELEGRAM_WEBHOOK_URL"))
	if webhookURL != "" {
		startWebhook(bot, baseURL, webhookURL)
	} else {
		startPolling(bot, baseURL)
	}
}
