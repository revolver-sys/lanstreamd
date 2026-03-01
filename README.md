# lanstreamd

LAN streaming helper daemon: forward proxy + observability.

## What it is
- HTTP forward proxy
- HTTPS tunneling via CONNECT (no TLS interception / no MITM)
- Local status API for health and counters

## Use cases
- Observe unstable network behavior during streaming
- See "bad minutes": connect failures, upstream errors, timeouts
- View status from any device in the local network

## Quick start
go mod tidy
go run ./cmd/lanstreamd -config configs/example.yaml

API:
GET /health
GET /status

curl http://127.0.0.1:8080/health
curl http://127.0.0.1:8080/status
Proxy check:
curl -x http://127.0.0.1:3128 https://api.ipify.org?format=text

Roadmap: see docs/roadmap.md
