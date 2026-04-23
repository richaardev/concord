package concord

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/omit"
)

type SlashCommand struct {
	Name                     string
	NameLocalizations        map[discord.Locale]string
	Description              string
	DescriptionLocalizations map[discord.Locale]string

	DefaultMemberPermissions omit.Omit[*discord.Permissions]
	NSFW                     *bool

	Annotations map[string]string

	PreRunE      func(event *events.ApplicationCommandInteractionCreate) bool
	RunE         func(event *events.ApplicationCommandInteractionCreate) error
	AutoComplete func(event *events.AutocompleteInteractionCreate) error

	Options         []discord.ApplicationCommandOption
	Contexts        []discord.InteractionContextType
	IntegrationType []discord.ApplicationIntegrationType

	SubCommands []SubCommand
}

type SubCommand struct {
	Name                     string                    `json:"name"`
	NameLocalizations        map[discord.Locale]string `json:"name_localizations,omitempty"`
	Description              string                    `json:"description"`
	DescriptionLocalizations map[discord.Locale]string `json:"description_localizations,omitempty"`

	Annotations map[string]string                                             `json:"-"`
	RunE        func(event *events.ApplicationCommandInteractionCreate) error `json:"-"`
	Options     []discord.ApplicationCommandOption                            `json:"options,omitempty"`
}

type MessageCommand struct {
	Name              string
	NameLocalizations map[discord.Locale]string

	DefaultMemberPermissions omit.Omit[*discord.Permissions]
	NSFW                     *bool

	IntegrationTypes []discord.ApplicationIntegrationType
	Contexts         []discord.InteractionContextType

	Annotations map[string]string

	PreRunE func(event *events.ApplicationCommandInteractionCreate) bool
	RunE    func(event *events.ApplicationCommandInteractionCreate) error
}

type MessageComponent struct {
	Annotations map[string]string
	RunE        func(event *events.ComponentInteractionCreate) error

	CustomID      string
	CustomIDRegex string
}

type ModalComponent struct {
	Annotations map[string]string
	RunE        func(event *events.ModalSubmitInteractionCreate) error

	CustomID      string
	CustomIDRegex string
}

func SlashCommandToCreate(slash *SlashCommand) discord.SlashCommandCreate {
	create := discord.SlashCommandCreate{
		Name:                     slash.Name,
		NameLocalizations:        slash.NameLocalizations,
		Description:              slash.Description,
		DescriptionLocalizations: slash.DescriptionLocalizations,
		Options:                  slash.Options,
		DefaultMemberPermissions: slash.DefaultMemberPermissions,
		IntegrationTypes:         slash.IntegrationType,
		Contexts:                 slash.Contexts,
		NSFW:                     slash.NSFW,
	}

	for _, subcommand := range slash.SubCommands {
		create.Options = append(create.Options, discord.ApplicationCommandOptionSubCommand{
			Name:        subcommand.Name,
			Description: subcommand.Description,
			Options:     subcommand.Options,
		})
	}

	return create
}

func SlashSubCommandToCreate(c *SubCommand) discord.ApplicationCommandOption {
	return discord.ApplicationCommandOptionSubCommand{
		Name:                     c.Name,
		NameLocalizations:        c.NameLocalizations,
		Description:              c.Description,
		DescriptionLocalizations: c.DescriptionLocalizations,
		Options:                  c.Options,
	}
}

func CommandListToCreate(c []*SlashCommand) []discord.SlashCommandCreate {
	var result []discord.SlashCommandCreate
	for _, command := range c {
		result = append(result, SlashCommandToCreate(command))
	}

	return result
}

func MessageCommandToCreate(command *MessageCommand) discord.MessageCommandCreate {
	return discord.MessageCommandCreate{
		Name:                     command.Name,
		NameLocalizations:        command.NameLocalizations,
		DefaultMemberPermissions: command.DefaultMemberPermissions,
		IntegrationTypes:         command.IntegrationTypes,
		Contexts:                 command.Contexts,
		NSFW:                     command.NSFW,
	}
}

func MessageCommandListToCreate(c []*MessageCommand) []discord.MessageCommandCreate {
	var result []discord.MessageCommandCreate
	for _, command := range c {
		result = append(result, MessageCommandToCreate(command))
	}

	return result
}
