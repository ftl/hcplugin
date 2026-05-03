package action

import (
	"fmt"

	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

const DoActionUUID = "com.thecodingflow.hcplugin.doaction"

func init() { Factories[DoActionUUID] = NewDoAction }

type DoAction struct {
	context string
	client  ClientAccessor
	deck    Deck
	status  StatusRecorder
}

func NewDoAction(context string, client ClientAccessor, deck Deck, status StatusRecorder) Action {
	return &DoAction{context: context, client: client, deck: deck, status: status}
}

func (a *DoAction) parseSettings(settings map[string]any) string {
	s, _ := settings["actionId"].(string)
	return s
}

func (a *DoAction) KeyDown(p *sdk.ReceivedEventPayload) error  { return a.fire(p) }
func (a *DoAction) DialDown(p *sdk.ReceivedEventPayload) error { return a.fire(p) }

func (a *DoAction) fire(p *sdk.ReceivedEventPayload) error {
	actionID := a.parseSettings(p.Settings)
	if actionID == "" {
		a.status.Record(a.deck, a.context, fmt.Errorf("actionId not configured"))
		return nil
	}
	a.status.Record(a.deck, a.context, a.client().Do(actionID))
	return nil
}
