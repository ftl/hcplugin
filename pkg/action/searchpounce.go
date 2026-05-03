package action

import (
	sdk "github.com/SkYNewZ/streamdeck-sdk"
)

const SearchPounceUUID = "com.thecodingflow.hcplugin.searchpounce"

func init() { Factories[SearchPounceUUID] = NewSearchPounce }

type SearchPounce struct {
	context string
	client  ClientAccessor
	deck    Deck
	status  StatusRecorder
}

func NewSearchPounce(context string, client ClientAccessor, deck Deck, status StatusRecorder) Action {
	return &SearchPounce{context: context, client: client, deck: deck, status: status}
}

func (a *SearchPounce) DialDown(_ *sdk.ReceivedEventPayload) error {
	a.status.Record(a.deck, a.context, a.client().Do("entry.next_esm_step"))
	return nil
}

func (a *SearchPounce) DialRotate(p *sdk.ReceivedEventPayload) error {
	ticks := p.Ticks
	if ticks == 0 {
		return nil
	}
	actionID := "bandmap.goto_next_spot_up"
	n := ticks
	if n < 0 {
		actionID = "bandmap.goto_next_spot_down"
		n = -n
	}
	var lastErr error
	if err := a.client().Do(actionID); err != nil {
		lastErr = err
	}
	a.status.Record(a.deck, a.context, lastErr)
	return nil
}
