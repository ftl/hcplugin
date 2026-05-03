package action

import (
	"fmt"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

const SendTextUUID = "com.thecodingflow.hcplugin.sendtext"

func init() { Factories[SendTextUUID] = NewSendText }

type SendText struct {
	context string
	client  ClientAccessor
	deck    Deck
	status  StatusRecorder
}

func NewSendText(context string, client ClientAccessor, deck Deck, status StatusRecorder) Action {
	return &SendText{context: context, client: client, deck: deck, status: status}
}

func (a *SendText) parseSettings(settings map[string]any) string {
	s, _ := settings["text"].(string)
	return s
}

func (a *SendText) KeyDown(p *sdk.ReceivedEventPayload) error  { return a.fire(p) }
func (a *SendText) DialDown(p *sdk.ReceivedEventPayload) error { return a.fire(p) }

func (a *SendText) fire(p *sdk.ReceivedEventPayload) error {
	text := a.parseSettings(p.Settings)
	if text == "" {
		a.status.Record(a.deck, a.context, fmt.Errorf("text not configured"))
		return nil
	}
	a.status.Record(a.deck, a.context, a.client().Send(text))
	return nil
}
