$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$readmePath = Join-Path $repoRoot "README.md"

$env:PORT = "30001"
$env:RAG_DOC_PATH = $readmePath

if (-not $env:CAPTCHA_SECRET) {
  $env:CAPTCHA_SECRET = "dev-secret"
}

Write-Host "Starting backend on port $env:PORT"
go run main.go
