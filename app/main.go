package main

import (
	"log"
	"os"
	"strconv"

	"github.com/celestix/gotgproto"
	"github.com/celestix/gotgproto/dispatcher/handlers"
	"github.com/celestix/gotgproto/dispatcher/handlers/filters"
	"github.com/celestix/gotgproto/ext"
	"github.com/celestix/gotgproto/sessionMaker"
	"github.com/glebarez/sqlite"
	"github.com/gotd/td/tg"
)

func main() {
	appIDStr := os.Getenv("TELEGRAM_APP_ID")
	appHash := os.Getenv("TELEGRAM_HASH_ID")
	appToken := os.Getenv("TELEGRAM_APP_TOKEN")

	if appIDStr == "" || appHash == "" || appToken == "" {
		// logger.Sugar().Fatalln("Environment variables TELEGRAM_APP_ID or TELEGRAM_HASH_ID or TELEGRAM_APP_TOKEN are not set")
		panic("Environment variables TELEGRAM_APP_ID or TELEGRAM_HASH_ID or TELEGRAM_APP_TOKEN are not set")
	}

	appID, appIDError := strconv.Atoi(appIDStr)
	if appIDError != nil {
		// logger.Sugar().Fatalf("Error converting TELEGRAM_APP_ID to integer: %v", appIDError)
		panic("Error converting TELEGRAM_APP_ID to integer")
	}

	client, err := gotgproto.NewClient(
		// Get AppID from https://my.telegram.org/apps
		appID,
		// Get ApiHash from https://my.telegram.org/apps
		appHash,
		// ClientType, as we defined above
		gotgproto.ClientTypeBot(appToken),
		// Optional parameters of client
		&gotgproto.ClientOpts{
			Session: sessionMaker.SqlSession(sqlite.Open("bot.session")),
		},
	)

	if err != nil {
		log.Fatalln("failed to start client:", err)
	}

	dispatcher := client.Dispatcher

	// This Message Handler will call our echo function on new messages
	dispatcher.AddHandlerToGroup(handlers.NewMessage(filters.Message.Text, echo), 1)

	client.Idle()

}

// Ref: https://github.com/celestix/gotgproto/blob/beta/examples/echo-bot/sqlite_session/main.go
func echo(ctx *ext.Context, update *ext.Update) error {
	msg := update.EffectiveMessage
	text := msg.Message.Message
	// userId := update.EffectiveMessage.ID
	chatId := update.EffectiveChat().GetID()

	var markupTable string = `
	<pre>
	| Tables   |      Are      |  Cool |
	|----------|:-------------:|------:|
	| col 1 is |  left-aligned | $1600 |
	| col 2 is |    centered   |   $12 |
	| col 3 is | right-aligned |    $1 |
	</pre>
	`

	_, err := ctx.Reply(update, text, &ext.ReplyOpts{})
	ctx.SendMessage(chatId, &tg.MessagesSendMessageRequest{
		Message: markupTable,
		// Peer: ... (No need of setting peer as we have passed chatId)
	})
	return err
}
