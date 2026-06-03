package pkg

import (
	"log"

	"poc-mini-loan/internal/config"

	"go.temporal.io/sdk/client"
)

func InitTemporalClient(cfg *config.Config) client.Client {
	log.Printf("🔌 Connecting to Temporal Server: %s", cfg.TemporalHost)

	c, err := client.Dial(client.Options{
		HostPort: cfg.TemporalHost,
	})

	if err != nil {
		log.Fatalf("❌ Cannot connect Temporal Server: %v", err)
	}

	log.Println("✅ Connected to Temporal Server successfully")
	return c
}
