package main

import (
	"app/commands"
	"log"
	"os"
	"strconv"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
	"go.uber.org/zap"
)

// https://github.com/celestix/gotgproto
func main() {

	logger, loggerError := zap.NewProduction() // or zap.NewDevelopment()
	if loggerError != nil {
		panic(loggerError)
	}
	defer logger.Sync() // Flush any buffered log entries

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

	client, clientError := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		appID,
		// Get ApiHash from https://my.telegram.org/apps
		appHash,
		// ClientType, as we defined above
		gotgproto.ClientTypeBot(appToken),
		// Optional parameters of client
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("bot.session")),
			Logger:  logger,
		},
	)

	if clientError != nil {
		log.Fatalln("failed to start client:", clientError)
	}

	dispatcher := client.Dispatcher

	// This Message Handler will call our echo function on new messages
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, commands.Echo), 1)

	client.Idle()

}
