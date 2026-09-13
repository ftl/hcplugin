package plugin

import (
	"log"
	"sync"
	"sync/atomic"

	sdk "github.com/SkYNewZ/streamdeck-sdk"

	"github.com/ftl/hcplugin/pkg/action"
	"github.com/ftl/hcplugin/pkg/remote"
)

// Deck is the subset of sdk.StreamDeck used by the plugin.
type Deck interface {
	action.Deck
	GetGlobalSettings(context string)
}

type Plugin struct {
	uuid          string
	deck          Deck
	client        atomic.Pointer[remote.Client]
	instances     map[string]any
	instancesLock sync.Mutex
	statusTracker *StatusTracker
}

func New(uuid string, deck Deck) *Plugin {
	log.Printf("started plugin %s", uuid)
	p := &Plugin{
		uuid:          uuid,
		deck:          deck,
		instances:     make(map[string]any),
		statusTracker: NewStatusTracker(),
	}
	p.client.Store(remote.New(DefaultPort))
	return p
}

func (p *Plugin) Start() {
	// GetGlobalSettings blocks until websocket is open, must run after Start() is called on the main goroutine
	go func() {
		p.deck.GetGlobalSettings(p.uuid)
	}()
}

func (p *Plugin) clientAccessor() *remote.Client {
	return p.client.Load()
}

func (p *Plugin) Handle(event *sdk.ReceivedEvent) error {
	log.Printf("a:%s c:%s e:%s:\n%+v\n\n", event.Action, event.Context, event.Event, event.Payload)

	if event.Event == sdk.DidReceiveGlobalSettings {
		globalSettings, err := parseGlobalSettings(event.Payload.Settings)
		if err != nil {
			return err
		}
		p.client.Store(remote.New(globalSettings.Port))
	}

	if event.Action == "" || event.Context == "" {
		return nil
	}

	if event.Event == sdk.WillDisappear {
		p.instancesLock.Lock()
		delete(p.instances, event.Context)
		p.instancesLock.Unlock()
		p.statusTracker.Forget(event.Context)
		return nil
	}

	p.instancesLock.Lock()
	inst, ok := p.instances[event.Context]
	if !ok {
		factory, ok2 := action.Factories[event.Action]
		if !ok2 {
			p.instancesLock.Unlock()
			log.Printf("unknown action %s", event.Action)
			return nil
		}
		inst = factory(event.Context, p.clientAccessor, p.deck, p.statusTracker)
		p.instances[event.Context] = inst
		log.Printf("NEW INSTANCE: %s = %T\n\n", event.Context, inst)
	}
	p.instancesLock.Unlock()

	switch event.Event {
	case sdk.KeyDown:
		return handle(inst, func(handler action.KeyDownHandler) error {
			return handler.KeyDown(event.Payload)
		})
	case sdk.KeyUp:
		return handle(inst, func(handler action.KeyUpHandler) error {
			return handler.KeyUp(event.Payload)
		})
	case sdk.DialRotate:
		return handle(inst, func(handler action.DialRotateHandler) error {
			return handler.DialRotate(event.Payload)
		})
	case sdk.DialDown:
		return handle(inst, func(handler action.DialDownHandler) error {
			return handler.DialDown(event.Payload)
		})
	default:
		return nil
	}
}

func handle[H any](handler any, call func(handler H) error) error {
	castedHandler, ok := handler.(H)
	if !ok {
		return nil
	}
	return call(castedHandler)
}
