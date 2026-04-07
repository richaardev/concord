package interactivity

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type InteractiveComponent interface {
	discord.InteractiveComponent
	Handler(e *events.ComponentInteractionCreate) error
}

// ============================================================================
// InteractiveModal
// ============================================================================

type InteractiveModal struct {
	discord.ModalCreate
	handler ModalHandler
}

func (i *InteractiveModal) Handler(e *events.ModalSubmitInteractionCreate) error {
	if i.handler == nil {
		return nil
	}

	return i.handler(e)
}

func (i *InteractiveModal) WithHandler(e func(e *events.ModalSubmitInteractionCreate) error) *InteractiveModal {
	i.handler = e
	return i
}

func (i *InteractiveModal) WithCustomID(customID string) *InteractiveModal {
	i.CustomID = customID
	return i
}

func (i *InteractiveModal) Create() discord.ModalCreate {
	return i.ModalCreate
}

// ============================================================================
// InteractiveButton
// ============================================================================

var (
	_ discord.Component            = (*InteractiveButton)(nil)
	_ discord.InteractiveComponent = (*InteractiveButton)(nil)
	_ InteractiveComponent         = (*InteractiveButton)(nil)
)

type InteractiveButton struct {
	discord.ButtonComponent
	handler ComponentHandler
}

func (i *InteractiveButton) Handler(e *events.ComponentInteractionCreate) error {
	if i.handler == nil {
		return nil
	}

	return i.handler(e)
}

func (i *InteractiveButton) WithHandler(e func(e *events.ComponentInteractionCreate) error) *InteractiveButton {
	i.handler = e
	return i
}

func (i *InteractiveButton) WithMessageUpdate(msg discord.MessageUpdate) *InteractiveButton {
	i.handler = func(e *events.ComponentInteractionCreate) error {
		return e.UpdateMessage(msg)
	}

	return i
}

func (i *InteractiveButton) WithID(id int) *InteractiveButton {
	i.ID = id
	return i
}

func (i *InteractiveButton) WithStyle(style discord.ButtonStyle) *InteractiveButton {
	i.Style = style
	return i
}

func (i *InteractiveButton) WithLabel(label string) *InteractiveButton {
	i.Label = label
	return i
}

func (i *InteractiveButton) WithEmoji(emoji discord.ComponentEmoji) *InteractiveButton {
	i.Emoji = &emoji
	return i
}

func (i *InteractiveButton) WithCustomID(customID string) *InteractiveButton {
	i.CustomID = customID
	return i
}

func (i *InteractiveButton) WithURL(url string) *InteractiveButton {
	i.URL = url
	return i
}

func (i *InteractiveButton) AsEnabled() *InteractiveButton {
	i.Disabled = false
	return i
}

func (i *InteractiveButton) AsDisabled() *InteractiveButton {
	i.Disabled = true
	return i
}

func (i *InteractiveButton) WithDisabled(disabled bool) *InteractiveButton {
	i.Disabled = disabled
	return i
}

// ============================================================================
// InteractiveStringSelectMenu
// ============================================================================

var (
	_ discord.Component            = (*InteractiveStringSelectMenu)(nil)
	_ discord.InteractiveComponent = (*InteractiveStringSelectMenu)(nil)
	_ InteractiveComponent         = (*InteractiveStringSelectMenu)(nil)
)

type InteractiveStringSelectMenu struct {
	discord.StringSelectMenuComponent
	handler ComponentHandler
}

func (i *InteractiveStringSelectMenu) Handler(e *events.ComponentInteractionCreate) error {
	if i.handler == nil {
		return nil
	}

	return i.handler(e)
}

func (i *InteractiveStringSelectMenu) WithHandler(fn func(e *events.ComponentInteractionCreate) error) *InteractiveStringSelectMenu {
	i.handler = fn
	return i
}

func (i *InteractiveStringSelectMenu) WithMaxValues(maxValues int) *InteractiveStringSelectMenu {
	i.MaxValues = maxValues
	return i
}

func (i *InteractiveStringSelectMenu) WithMinValues(minValues int) *InteractiveStringSelectMenu {
	i.MinValues = &minValues
	return i
}

func (i *InteractiveStringSelectMenu) WithRequired(required bool) *InteractiveStringSelectMenu {
	i.Required = required
	return i
}

func (i *InteractiveStringSelectMenu) WithDisabled(disabled bool) *InteractiveStringSelectMenu {
	i.Disabled = disabled
	return i
}

func (i *InteractiveStringSelectMenu) AsEnabled() *InteractiveStringSelectMenu {
	i.Disabled = false
	return i
}

func (i *InteractiveStringSelectMenu) AsDisabled() *InteractiveStringSelectMenu {
	i.Disabled = true
	return i
}

func (i *InteractiveStringSelectMenu) AddOptions(options ...discord.StringSelectMenuOption) *InteractiveStringSelectMenu {
	i.Options = append(i.Options, options...)
	return i
}

func (i *InteractiveStringSelectMenu) SetOptions(options ...discord.StringSelectMenuOption) *InteractiveStringSelectMenu {
	i.Options = options
	return i
}

func (i *InteractiveStringSelectMenu) WithValues(values []string) *InteractiveStringSelectMenu {
	i.Values = values
	return i
}

func (i *InteractiveStringSelectMenu) WithPlaceholder(placeholder string) *InteractiveStringSelectMenu {
	i.Placeholder = placeholder
	return i
}

// ============================================================================
// InteractiveChannelSelectMenu
// ============================================================================

var (
	_ discord.Component            = (*InteractiveChannelSelectMenu)(nil)
	_ discord.InteractiveComponent = (*InteractiveChannelSelectMenu)(nil)
	_ InteractiveComponent         = (*InteractiveChannelSelectMenu)(nil)
)

type InteractiveChannelSelectMenu struct {
	discord.ChannelSelectMenuComponent
	handler ComponentHandler
}

func (i *InteractiveChannelSelectMenu) Handler(e *events.ComponentInteractionCreate) error {
	if i.handler == nil {
		return nil
	}

	return i.handler(e)
}

func (i *InteractiveChannelSelectMenu) WithHandler(fn func(e *events.ComponentInteractionCreate) error) *InteractiveChannelSelectMenu {
	i.handler = fn
	return i
}

func (i *InteractiveChannelSelectMenu) WithChannelTypes(channelTypes ...discord.ChannelType) *InteractiveChannelSelectMenu {
	i.ChannelTypes = channelTypes
	return i
}

func (i *InteractiveChannelSelectMenu) WithMaxValues(maxValues int) *InteractiveChannelSelectMenu {
	i.MaxValues = maxValues
	return i
}

func (i *InteractiveChannelSelectMenu) WithMinValues(minValues int) *InteractiveChannelSelectMenu {
	i.MinValues = &minValues
	return i
}

func (i *InteractiveChannelSelectMenu) WithRequired(required bool) *InteractiveChannelSelectMenu {
	i.Required = required
	return i
}

func (i *InteractiveChannelSelectMenu) WithDisabled(disabled bool) *InteractiveChannelSelectMenu {
	i.Disabled = disabled
	return i
}

func (i *InteractiveChannelSelectMenu) AsEnabled() *InteractiveChannelSelectMenu {
	i.Disabled = false
	return i
}

func (i *InteractiveChannelSelectMenu) AsDisabled() *InteractiveChannelSelectMenu {
	i.Disabled = true
	return i
}

func (i *InteractiveChannelSelectMenu) SetDefaultValues(defaultValues ...snowflake.ID) *InteractiveChannelSelectMenu {
	values := make([]discord.SelectMenuDefaultValue, 0, len(defaultValues))
	for _, value := range defaultValues {
		values = append(values, discord.NewSelectMenuDefaultChannel(value))
	}
	i.DefaultValues = values
	return i
}

func (i *InteractiveChannelSelectMenu) AddDefaultValue(defaultValue snowflake.ID) *InteractiveChannelSelectMenu {
	i.DefaultValues = append(i.DefaultValues, discord.NewSelectMenuDefaultChannel(defaultValue))
	return i
}

func (i *InteractiveChannelSelectMenu) WithPlaceholder(placeholder string) *InteractiveChannelSelectMenu {
	i.Placeholder = placeholder
	return i
}

// ============================================================================
// InteractiveRoleSelectMenu
// ============================================================================

var (
	_ discord.Component            = (*InteractiveRoleSelectMenu)(nil)
	_ discord.InteractiveComponent = (*InteractiveRoleSelectMenu)(nil)
	_ InteractiveComponent         = (*InteractiveRoleSelectMenu)(nil)
)

type InteractiveRoleSelectMenu struct {
	discord.RoleSelectMenuComponent
	handler ComponentHandler
}

func (i *InteractiveRoleSelectMenu) Handler(e *events.ComponentInteractionCreate) error {
	if i.handler == nil {
		return nil
	}

	return i.handler(e)
}

func (i *InteractiveRoleSelectMenu) WithHandler(fn func(e *events.ComponentInteractionCreate) error) *InteractiveRoleSelectMenu {
	i.handler = fn
	return i
}

func (i *InteractiveRoleSelectMenu) WithMaxValues(maxValues int) *InteractiveRoleSelectMenu {
	i.MaxValues = maxValues
	return i
}

func (i *InteractiveRoleSelectMenu) WithMinValues(minValues int) *InteractiveRoleSelectMenu {
	i.MinValues = &minValues
	return i
}

func (i *InteractiveRoleSelectMenu) WithRequired(required bool) *InteractiveRoleSelectMenu {
	i.Required = required
	return i
}

func (i *InteractiveRoleSelectMenu) WithDisabled(disabled bool) *InteractiveRoleSelectMenu {
	i.Disabled = disabled
	return i
}

func (i *InteractiveRoleSelectMenu) AsEnabled() *InteractiveRoleSelectMenu {
	i.Disabled = false
	return i
}

func (i *InteractiveRoleSelectMenu) AsDisabled() *InteractiveRoleSelectMenu {
	i.Disabled = true
	return i
}

func (i *InteractiveRoleSelectMenu) SetDefaultValues(defaultValues ...snowflake.ID) *InteractiveRoleSelectMenu {
	values := make([]discord.SelectMenuDefaultValue, 0, len(defaultValues))
	for _, value := range defaultValues {
		values = append(values, discord.NewSelectMenuDefaultRole(value))
	}
	i.DefaultValues = values
	return i
}

func (i *InteractiveRoleSelectMenu) AddDefaultValue(defaultValue snowflake.ID) *InteractiveRoleSelectMenu {
	i.DefaultValues = append(i.DefaultValues, discord.NewSelectMenuDefaultRole(defaultValue))
	return i
}

func (i *InteractiveRoleSelectMenu) WithPlaceholder(placeholder string) *InteractiveRoleSelectMenu {
	i.Placeholder = placeholder
	return i
}

// ============================================================================
// InteractiveUserSelectMenu
// ============================================================================

var (
	_ discord.Component            = (*InteractiveUserSelectMenu)(nil)
	_ discord.InteractiveComponent = (*InteractiveUserSelectMenu)(nil)
	_ InteractiveComponent         = (*InteractiveUserSelectMenu)(nil)
)

type InteractiveUserSelectMenu struct {
	discord.UserSelectMenuComponent
	handler ComponentHandler
}

func (i *InteractiveUserSelectMenu) Handler(e *events.ComponentInteractionCreate) error {
	if i.handler == nil {
		return nil
	}

	return i.handler(e)
}

func (i *InteractiveUserSelectMenu) WithHandler(fn func(e *events.ComponentInteractionCreate) error) *InteractiveUserSelectMenu {
	i.handler = fn
	return i
}

func (i *InteractiveUserSelectMenu) WithMaxValues(maxValues int) *InteractiveUserSelectMenu {
	i.MaxValues = maxValues
	return i
}

func (i *InteractiveUserSelectMenu) WithMinValues(minValues int) *InteractiveUserSelectMenu {
	i.MinValues = &minValues
	return i
}

func (i *InteractiveUserSelectMenu) WithRequired(required bool) *InteractiveUserSelectMenu {
	i.Required = required
	return i
}

func (i *InteractiveUserSelectMenu) WithDisabled(disabled bool) *InteractiveUserSelectMenu {
	i.Disabled = disabled
	return i
}

func (i *InteractiveUserSelectMenu) AsEnabled() *InteractiveUserSelectMenu {
	i.Disabled = false
	return i
}

func (i *InteractiveUserSelectMenu) AsDisabled() *InteractiveUserSelectMenu {
	i.Disabled = true
	return i
}

func (i *InteractiveUserSelectMenu) SetDefaultValues(defaultValues ...snowflake.ID) *InteractiveUserSelectMenu {
	values := make([]discord.SelectMenuDefaultValue, 0, len(defaultValues))
	for _, value := range defaultValues {
		values = append(values, discord.NewSelectMenuDefaultUser(value))
	}
	i.DefaultValues = values
	return i
}

func (i *InteractiveUserSelectMenu) AddDefaultValue(defaultValue snowflake.ID) *InteractiveUserSelectMenu {
	i.DefaultValues = append(i.DefaultValues, discord.NewSelectMenuDefaultUser(defaultValue))
	return i
}

func (i *InteractiveUserSelectMenu) WithPlaceholder(placeholder string) *InteractiveUserSelectMenu {
	i.Placeholder = placeholder
	return i
}
