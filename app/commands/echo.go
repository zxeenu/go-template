package commands

import (
	"github.com/celestix/gotgproto/ext"
	"github.com/gotd/td/tg"
)

// Ref: https://github.com/celestix/gotgproto/blob/beta/examples/echo-bot/sqlite_session/main.go
func Echo(ctx *ext.Context, update *ext.Update) error {
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
