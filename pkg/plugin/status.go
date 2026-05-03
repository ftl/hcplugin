package plugin

import (
	"sync"

	"github.com/ftl/hcplugin/pkg/action"
)

type StatusTracker struct {
	mu    sync.Mutex
	state map[string]bool // context → last was OK (true) or error (false); absent = no prior call
}

func NewStatusTracker() *StatusTracker {
	return &StatusTracker{
		state: make(map[string]bool),
	}
}

func (t *StatusTracker) Record(deck action.Deck, ctx string, err error) {
	t.mu.Lock()
	prevOK, hasPrev := t.state[ctx]
	if err != nil {
		t.state[ctx] = false
	} else {
		t.state[ctx] = true
	}
	t.mu.Unlock()

	if err != nil {
		deck.Alert(ctx)
		return
	}
	if hasPrev && !prevOK {
		deck.ShowOK(ctx)
	}
}

func (t *StatusTracker) Forget(ctx string) {
	t.mu.Lock()
	delete(t.state, ctx)
	t.mu.Unlock()
}
