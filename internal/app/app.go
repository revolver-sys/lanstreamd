package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/revolver-sys/lanstreamd/internal/config"
	"github.com/revolver-sys/lanstreamd/internal/httpapi"
	"github.com/revolver-sys/lanstreamd/internal/metrics"
	"github.com/revolver-sys/lanstreamd/internal/proxy"
)

type App struct {
	Cfg *config.Config
	M   *metrics.Metrics
}

func New(cfg *config.Config) *App {
	return &App{
		Cfg: cfg,
		M:   metrics.New(),
	}
}

func (a *App) Run(ctx context.Context) error {
	tr := proxy.NewTransport(
		time.Duration(a.Cfg.Proxy.ConnectTimeoutMs)*time.Millisecond,
		time.Duration(a.Cfg.Proxy.IdleConnTimeoutMs)*time.Millisecond,
	)

	proxySrv := &http.Server{
		Addr: a.Cfg.Proxy.Listen,
		Handler: &proxy.Handler{
			Transport: tr,
			Metrics:   a.M,
		},
	}

	apiSrv := &http.Server{
		Addr:    a.Cfg.API.Listen,
		Handler: (&httpapi.Server{M: a.M}).Routes(),
	}

	errCh := make(chan error, 2)

	go func() { errCh <- proxySrv.ListenAndServe() }()
	go func() { errCh <- apiSrv.ListenAndServe() }()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = proxySrv.Shutdown(shutdownCtx)
		_ = apiSrv.Shutdown(shutdownCtx)
		return nil
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	}
}
