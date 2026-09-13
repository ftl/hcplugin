package action

import (
	"fmt"
	"time"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

const DoActionUUID = "com.thecodingflow.hcplugin.doaction"

func init() { Factories[DoActionUUID] = NewDoAction }

type DoAction struct {
	context string
	client  ClientAccessor
	deck    Deck
	status  StatusRecorder
	press   keyPress
}

func NewDoAction(context string, client ClientAccessor, deck Deck, status StatusRecorder) Action {
	return &DoAction{context: context, client: client, deck: deck, status: status}
}

func (a *DoAction) parseSettings(settings map[string]any, longPress bool) string {
	if longPress {
		// a long press without its own action ID does the same as a short press
		if s, _ := settings["longPressActionId"].(string); s != "" {
			return s
		}
	}
	s, _ := settings["actionId"].(string)
	return s
}

// the key up decides about the length of the press, therefore the action runs on the key up
func (a *DoAction) KeyDown(_ *sdk.ReceivedEventPayload) error {
	a.press.down(time.Now())
	return nil
}

func (a *DoAction) KeyUp(p *sdk.ReceivedEventPayload) error {
	return a.fire(p, a.press.up(time.Now()))
}

func (a *DoAction) DialDown(p *sdk.ReceivedEventPayload) error { return a.fire(p, false) }

func (a *DoAction) fire(p *sdk.ReceivedEventPayload, longPress bool) error {
	actionID := a.parseSettings(p.Settings, longPress)
	if actionID == "" {
		a.status.Record(a.deck, a.context, fmt.Errorf("actionId not configured"))
		return nil
	}
	a.status.Record(a.deck, a.context, a.client().Do(actionID))
	return nil
}
