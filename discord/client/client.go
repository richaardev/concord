package client

import (
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/richaardev/concord/database"
	"github.com/richaardev/concord/managers"
)

type BotClientOption func(*BotClient)

type BotClient struct {
	*disgobot.Client

	InteractionManager managers.InteractionManager
	InteractionWorker  managers.InteractionWorker
	ModuleManager      managers.ModuleManager

	Database database.LocalDatabase
}

func WithDatabase(db database.LocalDatabase) BotClientOption {
	return func(c *BotClient) {
		c.Database = db
	}
}

func NewClient(base *disgobot.Client, opts ...BotClientOption) *BotClient {
	interactionManager := managers.NewInteractionManager()
	interactionWorker := managers.NewInteractionWorker(interactionManager, 128)

	moduleManager := managers.NewModuleManager()

	c := &BotClient{
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
