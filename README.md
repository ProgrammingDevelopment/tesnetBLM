# War Tiket Engine

War Tiket Engine is a high-performance, real-time ticket queuing system designed to handle massive traffic spikes during "ticket wars." It includes a Go backend engine, a modern Next.js frontend, and integration guidance for external notification bots.

## Project Overview

Core features:
- High-performance queuing using Go concurrency and Redis atomic counters.
- Real-time status and quota updates in the frontend.
- Fairness guarantee with first-come, first-served logic enforced at the data layer.
- Scalable architecture suitable for microservices.
- Telegram bot integration guidance.

## Technology Stack

| Component | Technology | Description |
| --- | --- | --- |
| Backend | Go (Golang) | Core logic, API handling, concurrency management |
| Framework | Gin | Lightweight HTTP framework for Go |
| Database | PostgreSQL | Persistent storage for users and logs |
| Cache/Queue | Redis | In-memory atomic counters and blocking |
| Frontend | Next.js (React) | Server-side rendering and static site generation |
| Styling | Bootstrap 5 | Responsive UI components and layout |
| Containerization | Docker | Consistent local and deployment environments |

## Project Structure

```
tesnet/
|-- engine/              # Go backend
|   |-- handlers/        # API request handlers
|   |-- models/          # Data structures
|   |-- services/        # Business logic and DB connections
|   `-- main.go          # Entry point
|-- frontend/            # Next.js frontend
|   |-- src/             # Components and pages
|   |-- public/          # Static assets
|   `-- package.json     # Dependencies
|-- bot/                 # Telegram bot (optional integration)
|-- docker-compose.yml   # Infrastructure configuration
`-- README.md            # Project documentation
```

## Getting Started

Prerequisites:
- Docker and Docker Compose
- Go 1.22+
- Node.js 18+

1. Start infrastructure (Redis and Postgres).

```bash
docker-compose up -d
```

2. Start backend engine (Go).

```bash
cd engine
go mod tidy
go run main.go
```

Backend API: `http://localhost:8080`

3. Start frontend (Next.js).

```bash
cd frontend
npm install
npm run dev
```

Frontend app: `http://localhost:3000`

## Testing Guide

Manual testing via UI:
- Open `http://localhost:3000`.
- Use the Dev Controls at the bottom to simulate states (Pre-War, Loading, Success, Failure).
- Click "Ambil Antrean" to test the flow.
- Verify quota updates in real-time.

API testing via CLI:

```bash
curl -X POST http://localhost:8080/api/war \
  -H "Content-Type: application/json" \
  -d '{"user_id":"test_user_01","location_id":"graha-dipta","time_slot":"10:00"}'
```

## Telegram Bot Integration

This section explains how to add a Telegram bot to notify users of queue status and ticket success.

1. Create your bot in Telegram: open Telegram and search for `@BotFather`, send `/newbot`, name your bot, and save the API token.

2. Install the library.

```bash
go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5
```

3. Add the bot code (example in `engine/bot/main.go`).

```go
package main

import (
  "log"
  "os"

  tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
  bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
  if err != nil {
    log.Panic(err)
  }

  bot.Debug = true
  log.Printf("Authorized on account %s", bot.Self.UserName)

  u := tgbotapi.NewUpdate(0)
  u.Timeout = 60
  updates := bot.GetUpdatesChan(u)

  for update := range updates {
    if update.Message != nil && update.Message.IsCommand() {
      msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
      switch update.Message.Command() {
      case "start":
        msg.Text = "Welcome to War Tiket Bot!"
      case "status":
        msg.Text = "Backend is ONLINE. Quota available."
      default:
        msg.Text = "Unknown command."
      }
      bot.Send(msg)
    }
  }
}
```

4. Run the bot locally.

Windows PowerShell:

```powershell
$env:TELEGRAM_APITOKEN="YOUR_TOKEN_HERE"; go run bot/main.go
```

Linux or Mac:

```bash
export TELEGRAM_APITOKEN="YOUR_TOKEN_HERE" && go run bot/main.go
```

### Free Deployment Strategy

Option A: Oracle Cloud Always Free (recommended)
- Sign up for an Oracle Cloud Free Tier account.
- Create an Ampere Arm instance (Always Free).
- SSH into the instance.
- Transfer your Go binary or pull the repo.
- Run the bot using `systemd` or `tmux`.

Option B: Google Cloud Run (serverless via webhooks)
- Modify the bot to use webhooks instead of long polling.
- Build and push a Docker image to Google Container Registry.
- Deploy to Cloud Run.
- Set the webhook:

```bash
curl "https://api.telegram.org/bot<TOKEN>/setWebhook?url=<YOUR_CLOUD_RUN_URL>"
```

## System Evaluation (Penilaian Sistem War Tiket)

Executive summary:
Sistem War Tiket Engine adalah demonstrasi arsitektur modern untuk menangani beban trafik ekstrem dengan integritas data absolut. Kombinasi Go, Redis, dan PostgreSQL memastikan keadilan dan kecepatan sekaligus mencegah overselling.

Studi kasus dan permasalahan:
- Race Condition: Dua pengguna mengambil tiket terakhir bersamaan.
- Server Crash: Lonjakan trafik menyebabkan lock berlebihan.
- Bot and scripting abuse: Skrip otomatis memborong tiket.
- Inkonsistensi data: Stok terlihat ada di halaman depan, tetapi gagal saat pembayaran.

Solusi yang diterapkan:
- Atomic counting (Redis) dengan operasi `DECR` untuk mencegah race condition.
- Go concurrency (goroutines) untuk menangani puluhan ribu koneksi simultan.
- Strict validation layer di sisi server.
- Database transaction integrity (ACID) di PostgreSQL.

Mengapa sistem ini unggul:
- Performance: ribuan RPS tanpa degradasi signifikan.
- Reliability: zero overselling.
- Security: validasi server-side ketat terhadap manipulasi waktu dan input.
- Scalability: arsitektur stateless yang siap di-scale horizontal.

Pertanyaan kunci dewan juri:
"Bagaimana mendapatkan antrean sebelum default waktu dimulai dan terverifikasi serta tervalidasi?"

Jawaban teknis:
Secara desain sistem, hal tersebut tidak mungkin dilakukan pengguna umum karena penerapan Zero-Trust Time Validation.

Mekanisme keamanan:
- Server-side time authority: waktu antrean dikunci di backend dan diverifikasi dengan waktu server.
- Jika request masuk sebelum waktu mulai, server menolak dengan status 403 Forbidden atau 425 Too Early.
- Double-step validation: pre-check waktu, lalu pengecekan kuota atomik di Redis.
- Controlled bypass hanya untuk QA/Audit menggunakan whitelist token khusus.

Kesimpulan:
Validasi waktu server-side dan operasi atomik Redis memastikan tidak ada jalur belakang untuk mencuri start.

## Reference Links

External references:
- <https://antrean.logammulia.com/>
- <https://antrean.logammulia.com/register>
- <https://antrean.logammulia.com/login>

Technical references to cite in presentations:
- Redis atomic operations (Transactions) in Redis documentation.
- Go concurrency patterns (Pipelines and Cancellation) in the Go Blog.
- PostgreSQL ACID compliance (Transaction Isolation) in PostgreSQL documentation.
