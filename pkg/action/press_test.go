package action

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyPress(t *testing.T) {
	now := time.Date(2026, time.September, 13, 12, 0, 0, 0, time.UTC)
	tt := []struct {
		desc     string
		hold     time.Duration
		expected bool
	}{
		{desc: "short press", hold: 100 * time.Millisecond, expected: false},
		{desc: "just below the threshold", hold: LongPressThreshold - time.Millisecond, expected: false},
		{desc: "exactly the threshold", hold: LongPressThreshold, expected: true},
		{desc: "long press", hold: time.Second, expected: true},
	}
	for _, tc := range tt {
		t.Run(tc.desc, func(t *testing.T) {
			var press keyPress
			press.down(now)
			assert.Equal(t, tc.expected, press.up(now.Add(tc.hold)))
		})
	}
}

func TestKeyPress_UpWithoutDown(t *testing.T) {
	now := time.Date(2026, time.September, 13, 12, 0, 0, 0, time.UTC)
	var press keyPress

	assert.False(t, press.up(now), "a key up alone is a short press")

	press.down(now)
	assert.True(t, press.up(now.Add(time.Second)))
	assert.False(t, press.up(now.Add(time.Hour)), "the second key up finds no key down")
}
