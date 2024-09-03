package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/gotd/td/telegram"
)

const (
	appID    = 0
	appHash  = ""
	appToken = ""
)

func main() {

	// Create a Zap logger
	logger, err := zap.NewProduction() // or zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Flush any buffered log entries

	logger.Sugar().Infoln("Bot Initialized...")

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

		println(state)

		// It is only valid to use client while this function is not returned
		// and ctx is not cancelled.
		// api := client.API()

		// Now you can invoke MTProto RPC requests by calling the API.
		// ...

		// Return to close client connection and free up resources.
		return nil
	}); err != nil {
		println(err, "error")
		panic(err)
	}
	// Client is closed.
}
