package interactivity

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

var (
	interactivesMu sync.RWMutex
	interactives   = make(map[string]Interactivity)
)

var _ Interactivity = (*interactivityImpl)(nil)

type (
	ComponentHandlerMiddleware func(e *events.ComponentInteractionCreate) bool
	ComponentHandler           func(e *events.ComponentInteractionCreate) error
	ModalHandler               func(e *events.ModalSubmitInteractionCreate) error
)

type Interactivity interface {
	ID() string
	CustomID() string

	UserID() snowflake.ID
	ChannelID() snowflake.ID
	MessageID() snowflake.ID
	GuildID() snowflake.ID

	NewPrimaryButton(label string) *InteractiveButton
	NewSecondaryButton(label string) *InteractiveButton
	NewSuccessButton(label string) *InteractiveButton
	NewDangerButton(label string) *InteractiveButton

	NewStringSelectMenu(placeholder string, options ...discord.StringSelectMenuOption) *InteractiveStringSelectMenu
	NewChannelSelectMenu(placeholder string) *InteractiveChannelSelectMenu
	NewUserSelectMenu(placeholder string) *InteractiveUserSelectMenu
	NewRoleSelectMenu(placeholder string) *InteractiveRoleSelectMenu

	NewModal(title string, components ...discord.LayoutComponent) *InteractiveModal

	RemoveComponent(customID string)
	RemoveModal(customID string)
	Cancel()
}

type interactivityImpl struct {
	id       string
	customID string

	userID    snowflake.ID
	messageID snowflake.ID
	channelID snowflake.ID
	guildID   snowflake.ID

	client *bot.Client
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex

	resetIdle chan struct{}

	modalHandlers map[string]*InteractiveModal
	modalFilter   func(e *events.ModalSubmitInteractionCreate) bool

	middleware        ComponentHandlerMiddleware
	componentHandlers map[string]InteractiveComponent
	componentFilter   func(e *events.ComponentInteractionCreate) bool
}

func NewInteractivity(client *bot.Client, opts ...ManagerOption) Interactivity {
	id := randomString(64)
	config := ManagerOptions{}
	for _, opt := range opts {
		opt(&config)
	}

	var toCancel []Interactivity

	interactivesMu.RLock()
	if oldSession, exists := interactives[id]; exists {
		toCancel = append(toCancel, oldSession)
	}

	if config.CustomID != "" {
		for _, oldI := range interactives {
			if oldI.CustomID() == config.CustomID &&
				oldI.UserID() == config.UserID &&
				oldI.GuildID() == config.GuildID &&
				oldI.ChannelID() == config.ChannelID {
				toCancel = append(toCancel, oldI)
			}
		}
	}
	interactivesMu.RUnlock()

	for _, oldI := range toCancel {
		oldI.Cancel()
	}

	var ctx context.Context
	var ctxCancel context.CancelFunc

	if config.MaxTime > 0 {
		ctx, ctxCancel = context.WithTimeout(context.Background(), config.MaxTime)
	} else {
		ctx, ctxCancel = context.WithCancel(context.Background())
	}

	i := &interactivityImpl{
		id:       id,
		customID: config.CustomID,

		userID:    config.UserID,
		messageID: config.MessageID,
		channelID: config.ChannelID,
		guildID:   config.GuildID,

		client:            client,
		ctx:               ctx,
		cancel:            ctxCancel,
		componentHandlers: make(map[string]InteractiveComponent),
		modalHandlers:     make(map[string]*InteractiveModal),
	}

	i.cancel = func() {
		ctxCancel()

		interactivesMu.Lock()
		if current, ok := interactives[id]; ok && current == i {
			delete(interactives, id)
		}
		interactivesMu.Unlock()
	}

	if config.IdleTime > 0 {
		i.resetIdle = make(chan struct{}, 1)
		go i.monitorIdle(config.IdleTime)
	}

	go i.watchComponents()
	go i.watchModals()

	interactivesMu.Lock()
	interactives[id] = i
	interactivesMu.Unlock()

	return i
}

func (i *interactivityImpl) ID() string {
	return i.id
}

func (i *interactivityImpl) CustomID() string {
	return i.customID
}

func (i *interactivityImpl) ChannelID() snowflake.ID {
	return i.guildID
}

func (i *interactivityImpl) GuildID() snowflake.ID {
	return i.guildID
}

func (i *interactivityImpl) MessageID() snowflake.ID {
	return i.messageID
}

func (i *interactivityImpl) UserID() snowflake.ID {
	return i.userID
}

func (i *interactivityImpl) button(label string, style discord.ButtonStyle) *InteractiveButton {
	interactive := &InteractiveButton{
		ButtonComponent: discord.ButtonComponent{
			Label:    label,
			Style:    style,
			CustomID: "richaardev_loves_you_" + randomString(32),
		},
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.componentHandlers[interactive.GetCustomID()] = interactive

	return interactive
}

func (i *interactivityImpl) NewPrimaryButton(label string) *InteractiveButton {
	return i.button(label, discord.ButtonStylePrimary)
}

func (i *interactivityImpl) NewSecondaryButton(label string) *InteractiveButton {
	return i.button(label, discord.ButtonStyleSecondary)
}

func (i *interactivityImpl) NewSuccessButton(label string) *InteractiveButton {
	return i.button(label, discord.ButtonStyleSuccess)
}

func (i *interactivityImpl) NewDangerButton(label string) *InteractiveButton {
	return i.button(label, discord.ButtonStyleDanger)
}

func (i *interactivityImpl) NewStringSelectMenu(placeholder string, options ...discord.StringSelectMenuOption) *InteractiveStringSelectMenu {
	customID := "richaardev_loves_you_" + randomString(32)

	interactive := &InteractiveStringSelectMenu{
		StringSelectMenuComponent: discord.StringSelectMenuComponent{
			CustomID:    customID,
			Placeholder: placeholder,
			Options:     options,
		},
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.componentHandlers[interactive.GetCustomID()] = interactive

	return interactive
}

func (i *interactivityImpl) NewChannelSelectMenu(placeholder string) *InteractiveChannelSelectMenu {
	customID := "richaardev_loves_you_" + randomString(32)

	interactive := &InteractiveChannelSelectMenu{
		ChannelSelectMenuComponent: discord.ChannelSelectMenuComponent{
			CustomID:    customID,
			Placeholder: placeholder,
		},
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.componentHandlers[interactive.GetCustomID()] = interactive

	return interactive
}

func (i *interactivityImpl) NewRoleSelectMenu(placeholder string) *InteractiveRoleSelectMenu {
	customID := "richaardev_loves_you_" + randomString(32)

	interactive := &InteractiveRoleSelectMenu{
		RoleSelectMenuComponent: discord.RoleSelectMenuComponent{
			CustomID:    customID,
			Placeholder: placeholder,
		},
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.componentHandlers[interactive.GetCustomID()] = interactive

	return interactive
}

func (i *interactivityImpl) NewUserSelectMenu(placeholder string) *InteractiveUserSelectMenu {
	customID := "richaardev_loves_you_" + randomString(32)

	interactive := &InteractiveUserSelectMenu{
		UserSelectMenuComponent: discord.UserSelectMenuComponent{
			CustomID:    customID,
			Placeholder: placeholder,
		},
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.componentHandlers[interactive.GetCustomID()] = interactive

	return interactive
}

func (i *interactivityImpl) NewModal(title string, components ...discord.LayoutComponent) *InteractiveModal {
	interactive := &InteractiveModal{
		ModalCreate: discord.ModalCreate{
			Title:      title,
			Components: components,
			CustomID:   "richaardev_loves_you_" + randomString(32),
		},
	}

	i.mu.Lock()
	defer i.mu.Unlock()
	i.modalHandlers[interactive.CustomID] = interactive

	return interactive
}

func (i *interactivityImpl) RemoveComponent(customID string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.componentHandlers, customID)
}

func (i *interactivityImpl) RemoveModal(customID string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	delete(i.modalHandlers, customID)
}

func (i *interactivityImpl) Cancel() {
	if i.cancel != nil {
		i.cancel()
	}
}

func (i *interactivityImpl) notifyActivity() {
	if i.resetIdle != nil {
		select {
		case i.resetIdle <- struct{}{}:
		default:
		}
	}
}

func (i *interactivityImpl) monitorIdle(idleTime time.Duration) {
	timer := time.NewTimer(idleTime)
	defer timer.Stop()

	for {
		select {
		case <-i.ctx.Done():
			return
		case <-timer.C:
			i.cancel()
			return
		case <-i.resetIdle:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(idleTime)
		}
	}
}

func (i *interactivityImpl) watchComponents() {
	ch, cancel := bot.NewEventCollector(i.client, func(e *events.ComponentInteractionCreate) bool {
		if i.componentFilter != nil {
			return i.componentFilter(e)
		}
		return true
	})
	defer cancel()

	for {
		select {
		case <-i.ctx.Done():
			return
		case e := <-ch:
			i.mu.RLock()
			handler, ok := i.componentHandlers[e.Data.CustomID()]
			i.mu.RUnlock()

			if ok && e.Data.Type() == handler.Type() {
				i.notifyActivity()

				go func() {
					if i.middleware == nil || i.middleware(e) {
						if err := handler.Handler(e); err != nil {
							slog.Error("erro ao tentar executar o interactive", "error", err)
						}
					}
				}()
			}

		}
	}
}

func (i *interactivityImpl) watchModals() {
	ch, cancel := bot.NewEventCollector(i.client, func(e *events.ModalSubmitInteractionCreate) bool {
		if i.modalFilter != nil {
			return i.modalFilter(e)
		}
		return true
	})
	defer cancel()

	for {
		select {
		case <-i.ctx.Done():
			return
		case e := <-ch:
			i.mu.RLock()
			handler, ok := i.modalHandlers[e.Data.CustomID]
			i.mu.RUnlock()

			if ok {
				i.notifyActivity()
				go func() {
					if err := handler.Handler(e); err != nil {
						slog.Error("erro ao tentar executar o interactive", "error", err)
					}
				}()
			}
		}
	}
}
