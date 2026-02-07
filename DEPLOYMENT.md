# Deployment

## GitHub Repo
1. Initialize git and commit changes.
2. Push to a new GitHub repository.

## Cloudflare Pages (Frontend)
- Build command: `npx next build`
- Output directory: `out`
- Root directory: `frontend`
- Environment variable:
  - `NEXT_PUBLIC_API_BASE` = your backend URL (example: `https://api.yourdomain.com`)

## Backend
Run the Go backend on your server or VM (not on Cloudflare Pages).
- Set `PORT=30001`
- Set `RAG_DOC_PATH` to the README path
- Set `CAPTCHA_SECRET`

## Telegram Bot
- Set `TELEGRAM_APITOKEN`
- Set `BACKEND_BASE_URL` to your backend URL