package main

import (
	"fmt"
	"os"

	"codedock/cmd/codedock/client"
	"codedock/cmd/codedock/config"
)

func getClient() *client.Client {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	if cfg.ServerURL == "" || cfg.Token == "" {
		fmt.Println("Error: Not authenticated. Please run 'codedock login' first.")
		os.Exit(1)
	}
	return client.NewClient(cfg.ServerURL, cfg.Token)
}
