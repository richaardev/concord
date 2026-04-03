package interactivity

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/richaardev/concord/discordutils"
)

type PanelHandler func(e *events.ComponentInteractionCreate) error

type interactivityOptions struct {
	middleware ComponentHandlerMiddleware

	componentHandler ComponentHandler
	modalHandler     ModalHandler
}

type InteractivityOption func(*interactivityOptions)

func WithComponentHandler(handler ComponentHandler) InteractivityOption {
	return func(o *interactivityOptions) {
		o.componentHandler = handler
	}
}

func WithPanel(handler func() (discord.MessageCreate, error)) PanelHandler {
	return func(e *events.ComponentInteractionCreate) error {
		msg, err := handler()
		if err != nil {
			return fmt.Errorf("WithPanel: failed to generate panel message: %w", err)
		}

		return e.UpdateMessage(discordutils.ToMessageUpdate(msg))
	}
}

func WithMiddleware(handler ComponentHandlerMiddleware) InteractivityOption {
	return func(o *interactivityOptions) {
		o.middleware = handler
	}
}

func WithModalHandler(handler ModalHandler) InteractivityOption {
	return func(o *interactivityOptions) {
		o.modalHandler = handler
	}
}

func applyOptions(opts []InteractivityOption) interactivityOptions {
	options := interactivityOptions{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}

type ManagerOptions struct {
	CustomID string

	UserID    snowflake.ID
	GuildID   snowflake.ID
	ChannelID snowflake.ID
	MessageID snowflake.ID

	Middleware ComponentHandlerMiddleware

	MaxTime  time.Duration
	IdleTime time.Duration
}

type ManagerOption func(*ManagerOptions)

func WithOptions(opts ManagerOptions) ManagerOption {
	return func(o *ManagerOptions) {
		if opts.CustomID != "" {
			o.CustomID = opts.CustomID
		}

		if opts.ChannelID != 0 {
			o.ChannelID = opts.ChannelID
		}

		if opts.UserID != 0 {
			o.UserID = opts.UserID
		}

		if opts.MessageID != 0 {
			o.MessageID = opts.MessageID
		}

		if opts.MaxTime > 0 {
			o.MaxTime = opts.MaxTime
		}

		if opts.IdleTime > 0 {
			o.IdleTime = opts.IdleTime
		}

		if opts.Middleware != nil {
			o.Middleware = opts.Middleware
		}
	}
}
