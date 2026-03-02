package httpapi

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/revolver-sys/lanstreamd/internal/metrics"
	"github.com/revolver-sys/lanstreamd/internal/version"
)

type Server struct {
	M *metrics.Metrics
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		uptime := time.Since(s.M.StartTime).Round(time.Second).String()

		out := map[string]any{
			"service": "lanstreamd",
			"version": version.Version,
			"commit":  version.Commit,
			"uptime":  uptime,

			"requests_total":  s.M.RequestsTotal.Load(),
			"connect_total":   s.M.ConnectTotal.Load(),
			"connect_errors":  s.M.ConnectErrors.Load(),
			"upstream_errors": s.M.UpstreamErrors.Load(),
			"bytes_up":        s.M.BytesUp.Load(),
			"bytes_down":      s.M.BytesDown.Load(),
			"last_rtt_ms":     s.M.LastUpstreamRTTMs.Load(),
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	})

	return mux
}
