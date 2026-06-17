package main

import (
	"flag"
	"log"

	"github.com/shuyou-ai/shuyou-go/internal/app"
	"github.com/shuyou-ai/shuyou-go/internal/config"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("init app failed: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
