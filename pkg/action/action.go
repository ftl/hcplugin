package action

import (
	sdk "github.com/SkYNewZ/streamdeck-sdk"

	"github.com/ftl/hcplugin/pkg/remote"
)

type Action any

type ClientAccessor func() *remote.Client

type StatusRecorder interface {
	Record(deck Deck, ctx string, err error)
}

type Factory func(context string, client ClientAccessor, deck Deck, status StatusRecorder) Action

var Factories = map[string]Factory{}

type Deck interface {
	Alert(context string)
	ShowOK(context string)
}

type KeyDownHandler interface {
	KeyDown(*sdk.ReceivedEventPayload) error
}

type KeyUpHandler interface {
	KeyUp(*sdk.ReceivedEventPayload) error
}

type DialDownHandler interface {
	DialDown(*sdk.ReceivedEventPayload) error
}

type DialRotateHandler interface {
	DialRotate(*sdk.ReceivedEventPayload) error
}
