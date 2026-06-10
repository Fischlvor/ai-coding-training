package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	appconfig "ai-coding-training/internal/infrastructure/config"
	httpui "ai-coding-training/internal/interface/http"
)

func main() {
	cfgPath := appconfig.ConfigFile()
	cfg, err := appconfig.Load(cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	server, err := httpui.NewServer(cfg.HTTP.Port)
	if err != nil {
		log.Fatalf("create http server: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Printf("config-center is starting on %s\n", server.Addr())
	if err := server.Start(ctx); err != nil {
		log.Fatalf("run http server: %v", err)
	}
}
