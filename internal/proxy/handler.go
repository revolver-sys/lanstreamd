package proxy

import (
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/revolver-sys/lanstreamd/internal/metrics"
)

type Handler struct {
	Transport *http.Transport
	Metrics   *metrics.Metrics
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Metrics.RequestsTotal.Add(1)

	if r.Method == http.MethodConnect {
		h.handleConnect(w, r)
		return
	}

	h.handleHTTP(w, r)
}

func (h *Handler) handleHTTP(w http.ResponseWriter, r *http.Request) {
	// For http.Transport.RoundTrip in server handlers:
	// RequestURI must be empty.
	r2 := r.Clone(r.Context())
	r2.RequestURI = ""

	// Drop hop-by-hop headers (minimal, but important).
	stripHopByHopHeaders(r2.Header)

	start := time.Now()
	resp, err := h.Transport.RoundTrip(r2)
	if err != nil {
		h.Metrics.UpstreamErrors.Add(1)
		http.Error(w, "upstream error", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	stripHopByHopHeaders(resp.Header)

	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)

	n, _ := io.Copy(w, resp.Body)
	h.Metrics.BytesDown.Add(uint64(n))
	h.Metrics.NoteRTT(time.Since(start))
}

func (h *Handler) handleConnect(w http.ResponseWriter, r *http.Request) {
	h.Metrics.ConnectTotal.Add(1)

	// CONNECT host:port
	targetConn, err := net.DialTimeout("tcp", r.Host, 8*time.Second)
	if err != nil {
		h.Metrics.ConnectErrors.Add(1)
		http.Error(w, "connect failed", http.StatusBadGateway)
		return
	}

	hj, ok := w.(http.Hijacker)
	if !ok {
		_ = targetConn.Close()
		http.Error(w, "hijack not supported", http.StatusInternalServerError)
		return
	}

	clientConn, buf, err := hj.Hijack()
	if err != nil {
		_ = targetConn.Close()
		http.Error(w, "hijack failed", http.StatusInternalServerError)
		return
	}

	// Flush any buffered bytes in bufio.Writer returned by Hijack.
	if buf != nil {
		_ = buf.Flush()
	}

	// Tell client tunnel is established.
	_, _ = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	// Bidirectional copy with one-time close.
	var once sync.Once
	closeBoth := func() {
		_ = clientConn.Close()
		_ = targetConn.Close()
	}

	go func() {
		n, _ := io.Copy(targetConn, clientConn) // up
		h.Metrics.BytesUp.Add(uint64(n))
		once.Do(closeBoth)
	}()

	n, _ := io.Copy(clientConn, targetConn) // down
	h.Metrics.BytesDown.Add(uint64(n))
	once.Do(closeBoth)
}

func stripHopByHopHeaders(hdr http.Header) {
	// Standard hop-by-hop
	hdr.Del("Connection")
	hdr.Del("Proxy-Connection")
	hdr.Del("Keep-Alive")
	hdr.Del("Proxy-Authenticate")
	hdr.Del("Proxy-Authorization")
	hdr.Del("TE")
	hdr.Del("Trailers")
	hdr.Del("Transfer-Encoding")
	hdr.Del("Upgrade")

	// Also remove headers listed in Connection: ...
	if c := hdr.Get("Connection"); c != "" {
		for _, f := range strings.Split(c, ",") {
			if ff := strings.TrimSpace(f); ff != "" {
				hdr.Del(ff)
			}
		}
	}
}
