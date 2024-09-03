package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
)

func main() {

	// Create a Zap logger
	logger, err := zap.NewProduction() // or zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Flush any buffered log entries

	logger.Sugar().Infoln("Bot Initialized...")

	appIDStr := os.Getenv("TELEGRAM_APP_ID")
	appHash := os.Getenv("TELEGRAM_HASH_ID")
	appToken := os.Getenv("TELEGRAM_APP_TOKEN")

	if appIDStr == "" || appHash == "" || appToken == "" {
		logger.Sugar().Fatalln("Environment variables TELEGRAM_APP_ID or TELEGRAM_HASH_ID or TELEGRAM_APP_TOKEN are not set")
	}

	appID, appIDError := strconv.Atoi(appIDStr)
	if appIDError != nil {
		logger.Sugar().Fatalf("Error converting TELEGRAM_APP_ID to integer: %v", appIDError)
	}

	// https://core.telegram.org/api/obtaining_api_id
	client := telegram.NewClient(appID, appHash, telegram.Options{
		Logger: logger,
	})

	if err := client.Run(context.Background(), func(ctx context.Context) error {

		if _, err := client.Auth().Bot(ctx, appToken); err != nil {
			return err
		}
		state, err := client.API().UpdatesGetState(ctx)

		if err != nil {
			return err
		}

		// Printing the state properly
		fmt.Printf("%+v\n", state)

		// It is only valid to use client while this function is not returned
		// and ctx is not cancelled.
		// api := client.API()

		// Now you can invoke MTProto RPC requests by calling the API.
		// ...

		// Return to close client connection and free up resources.
		return nil
	}); err != nil {
		logger.Sugar().Fatalln(err)
	}
	// Client is closed.
}
