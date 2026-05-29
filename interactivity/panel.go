package interactivity

import (
	"context"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Message struct {
	Content         *string                     `json:"content,omitempty"`
	Embeds          *[]discord.Embed            `json:"embeds,omitempty"`
	Components      *[]discord.LayoutComponent  `json:"components,omitempty"`
	Attachments     *[]discord.AttachmentUpdate `json:"attachments,omitempty"`
	Files           []*discord.File             `json:"-"`
	AllowedMentions *discord.AllowedMentions    `json:"allowed_mentions,omitempty"`
	// Flags are the MessageFlags of the message.
	// Be careful not to override the current flags when editing messages from other users - this will result in a permission error.
	// Use MessageFlags.Add for flags like discord.MessageFlagSuppressEmbeds.
	Flags *discord.MessageFlags `json:"flags,omitempty"`

	TTS        bool           `json:"tts,omitempty"`
	StickerIDs []snowflake.ID `json:"sticker_ids,omitempty"`
}

type Cmd func() any

func ModalCmd(modal discord.ModalCreate) Cmd {
	return func() any {
		return modal
	}
}

type PanelModel interface {
	Init()

	// When returns nil, nil -> no changes will be made
	Update(e any) (PanelModel, Cmd)
	View() Message
}

func NewPanel(client *bot.Client, p PanelModel, opts ...PanelOption) PanelManager {
	var ctx context.Context
	var ctxCancel context.CancelFunc

	config := PanelOptions{}
	for _, opt := range opts {
		opt(&config)
	}

	if config.MaxTime > 0 {
		ctx, ctxCancel = context.WithTimeout(context.Background(), config.MaxTime)
	} else {
		ctx, ctxCancel = context.WithCancel(context.Background())
	}

	manager := &panelManagerImpl{
		client:  client,
		panelID: randomString(8),
		panel:   p,

		ctx:    ctx,
		cancel: ctxCancel,
	}

	manager.cancel = func() {
		ctxCancel()
	}

	if config.IdleTime > 0 {
		manager.resetIdle = make(chan struct{}, 1)
		go manager.monitorIdle(config.IdleTime)
	}

	go p.Init()
	go manager.watchComponents()
	go manager.watchModals()

	return manager
}

type PanelOptions struct {
	CustomID string

	UserID    snowflake.ID
	GuildID   snowflake.ID
	MessageID snowflake.ID

	// Maximum time the panel should be active.
	MaxTime time.Duration

	// Time after which the panel should be considered idle.
	IdleTime time.Duration
}

type PanelOption func(*PanelOptions)

func WithOptions(opts PanelOptions) PanelOption {
	return func(o *PanelOptions) {
		if opts.CustomID != "" {
			o.CustomID = opts.CustomID
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
	}
}
