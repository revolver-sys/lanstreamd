package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/revolver-sys/lanstreamd/internal/app"
	"github.com/revolver-sys/lanstreamd/internal/config"
)

func main() {
	cfgPath := flag.String("config", "configs/example.yaml", "path to config yaml")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Println("config error:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigCh)

	go func() {
		<-sigCh
		cancel()
	}()

	a := app.New(cfg)
	if err := a.Run(ctx); err != nil {
		fmt.Println("run error:", err)
		os.Exit(1)
	}
}
