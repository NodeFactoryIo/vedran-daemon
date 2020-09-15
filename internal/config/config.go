package config

import (
	"fmt"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

func InitMainConfig() {
	if os.Getenv("ENV") != "test" {
		setupSentry()
	}
}

func setupSentry() {
	dsn := os.Getenv("SENTRY_DSN")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:   dsn,
		Debug: false,
	})
	if err != nil {
		fmt.Println("Sentry init failed")
	}

	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
}
