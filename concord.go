package client

import (
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/richaardev/concord/managers"
)

type ConcordClientOption func(*ConcordClient)

type ConcordClient struct {
	*disgobot.Client

	InteractionManager managers.InteractionManager
	InteractionWorker  managers.InteractionWorker
	ModuleManager      managers.ModuleManager
}

func NewCordClient(base *disgobot.Client, opts ...ConcordClientOption) *ConcordClient {
	interactionManager := managers.NewInteractionManager()
	interactionWorker := managers.NewInteractionWorker(interactionManager, 128)

	moduleManager := managers.NewModuleManager()

	c := &ConcordClient{
		Client:             base,
		InteractionManager: interactionManager,
		InteractionWorker:  interactionWorker,
		ModuleManager:      moduleManager,
	}

	c.AddEventListeners(interactionWorker.Listeners()...)

	for _, opt := range opts {
		opt(c)
	}

	return c
}
