package client

import (
	"github.com/disgoorg/disgo"
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/richaardev/concord/managers"
)

type ConcordClient struct {
	*disgobot.Client

	InteractionManager managers.InteractionManager
	InteractionWorker  managers.InteractionWorker
	ModuleManager      managers.ModuleManager
}

func NewCordClient(token string, opts ...ConcordClientOption) (*ConcordClient, error) {
	options := applyConcordOptions(opts)

	disgoclient, err := disgo.New(token, options.disgoOpts...)
	if err != nil {
		return nil, err
	}

	interactionManager := managers.NewInteractionManager()
	interactionWorker := managers.NewInteractionWorker(interactionManager, 128)

	moduleManager := managers.NewModuleManager()

	c := &ConcordClient{
		Client:             disgoclient,
		InteractionManager: interactionManager,
		InteractionWorker:  interactionWorker,
		ModuleManager:      moduleManager,
	}

	c.AddEventListeners(interactionWorker.Listeners()...)

	return c, nil
}
