package action

import (
	"sync"
	"time"
)

const LongPressThreshold = 800 * time.Millisecond

// keyPress tells a long press from a short one by the time between key down and
// key up. The SDK dispatches every event in its own goroutine, therefore the
// state needs a lock.
type keyPress struct {
	lock   sync.Mutex
	downAt time.Time
}

func (p *keyPress) down(now time.Time) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.downAt = now
}

func (p *keyPress) up(now time.Time) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	downAt := p.downAt
	p.downAt = time.Time{}
	// a key up without a key down counts as a short press
	if downAt.IsZero() {
		return false
	}
	return now.Sub(downAt) >= LongPressThreshold
}
