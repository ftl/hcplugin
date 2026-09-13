package action

import (
	"fmt"
	"time"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

const SendTextUUID = "com.thecodingflow.hcplugin.sendtext"

func init() { Factories[SendTextUUID] = NewSendText }

type SendText struct {
	context string
	client  ClientAccessor
	deck    Deck
	status  StatusRecorder
	press   keyPress
}

func NewSendText(context string, client ClientAccessor, deck Deck, status StatusRecorder) Action {
	return &SendText{context: context, client: client, deck: deck, status: status}
}

func (a *SendText) parseSettings(settings map[string]any, longPress bool) string {
	if longPress {
		// a long press without its own text sends the same text as a short press
		if s, _ := settings["longPressText"].(string); s != "" {
			return s
		}
	}
	s, _ := settings["text"].(string)
	return s
}

// the key up decides about the length of the press, therefore the text goes out on the key up
func (a *SendText) KeyDown(_ *sdk.ReceivedEventPayload) error {
	a.press.down(time.Now())
	return nil
}

func (a *SendText) KeyUp(p *sdk.ReceivedEventPayload) error {
	return a.fire(p, a.press.up(time.Now()))
}

func (a *SendText) DialDown(p *sdk.ReceivedEventPayload) error { return a.fire(p, false) }

func (a *SendText) fire(p *sdk.ReceivedEventPayload, longPress bool) error {
	text := a.parseSettings(p.Settings, longPress)
	if text == "" {
		a.status.Record(a.deck, a.context, fmt.Errorf("text not configured"))
		return nil
	}
	a.status.Record(a.deck, a.context, a.client().Send(text))
	return nil
}
