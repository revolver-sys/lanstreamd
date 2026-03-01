# Roadmap

## v0 (current)
- Forward proxy (HTTP)
- HTTPS tunneling via CONNECT (no TLS interception)
- Status API (/health, /status)
- Basic counters + last RTT ms

## v0.1 Observability
- Per-host counters (top destinations)
- Better RTT metrics (connect time + first-byte latency)
- Rolling window "badness score" (errors/timeouts/connect failures)

## v0.2 Correlation
- Manual "dropout marker" endpoint (timestamped)
- Network probes scheduler (ping/tcp handshake/dns latency)
- Correlate markers with proxy error bursts and probes

## Later (separate project / integration)
- Traffic prioritization/routing should live in vpn-router-daemon:
  policy routing, pf anchors, per-profile routing modes.
