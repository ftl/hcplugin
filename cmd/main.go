package main

import (
	"log"

	sdk "github.com/SkYNewZ/streamdeck-sdk"

	"github.com/ftl/hcplugin/pkg/plugin"
)

var version = "development"

func main() {
	streamDeck, err := sdk.New()
	if err != nil {
		log.Fatal(err)
	}
	p := plugin.New(streamDeck.UUID, streamDeck)
	streamDeck.Handler(p)
	p.Start()
	streamDeck.Start()
}
