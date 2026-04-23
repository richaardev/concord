package concord

import (
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

var _ InteractionWorker = (*interactionWorkerImpl)(nil)

type InteractionWorker interface {
	GetFreeWorkers() int
	Worker(id int)

	Listeners() []bot.EventListener
}
type interactionWorkerImpl struct {
	manager                   InteractionManager
	interactionWorkersSize    int
	interactionWorkers        []bool
	InteractionHandlerChannel chan discord.Interaction
}

func NewInteractionWorker(manager InteractionManager, workersSize int) InteractionWorker {
	worker := &interactionWorkerImpl{
		manager:                manager,
		interactionWorkersSize: workersSize,
		interactionWorkers:     make([]bool, workersSize),

		InteractionHandlerChannel: make(chan discord.Interaction),
	}

	for i := range workersSize {
		go worker.Worker(i)
	}

	return worker
}

func (i *interactionWorkerImpl) GetFreeWorkers() int {
	freeWorkers := 0
	for _, worker := range i.interactionWorkers {
		if !worker {
			freeWorkers++
		}
	}
	return freeWorkers
}

func (i *interactionWorkerImpl) Worker(id int) {
	for interaction := range i.InteractionHandlerChannel {
		i.interactionWorkers[id] = true

		slog.Debug(
			"Running a new interaction",
			"Interaction Type", interaction.Type(),
			"Interaction ID", interaction.ID().String(),
			"User ID", interaction.User().ID.String(),
			"Free Workers", i.GetFreeWorkers(),
		)

		switch inte := interaction.(type) {
		case *events.AutocompleteInteractionCreate:
			i.handleAutocompleteInteraction(inte)
		case *events.ApplicationCommandInteractionCreate:
			i.handleCommandInteraction(inte)
		case *events.ComponentInteractionCreate:
			i.handleComponentInteration(inte)
		case *events.ModalSubmitInteractionCreate:
			i.handleModalSubmit(inte)
		default:
		}

		i.interactionWorkers[id] = false
	}
}

func (i *interactionWorkerImpl) handleCommandInteraction(event *events.ApplicationCommandInteractionCreate) {
	switch event.Data.(type) {
	case discord.SlashCommandInteractionData:
		i.handleSlashCommandInteraction(event)
	case discord.MessageCommandInteractionData:
		i.handleMessageCommandInteraction(event)
	default:
		slog.Warn("Unsupported application command interaction", slog.Any("Type", event.Data.Type()))
	}
}

func (i *interactionWorkerImpl) handleSlashCommandInteraction(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	var result *SlashCommand
	for _, c := range i.manager.ListSlashCommands() {
		if data.CommandName() == c.Name {
			result = c
			break
		}
	}

	if result != nil && result.PreRunE != nil {
		if doContinue := result.PreRunE(event); !doContinue {
			return
		}
	}

	if event.SlashCommandInteractionData().SubCommandName != nil {
		index := slices.IndexFunc(result.SubCommands, func(s SubCommand) bool {
			return s.Name == *event.SlashCommandInteractionData().SubCommandName
		})
		if index != -1 {
			subcommand := result.SubCommands[index]
			if err := subcommand.RunE(event); err != nil {
				slog.Error(err.Error())
			}
		}

		return
	}

	if result != nil {
		if err := result.RunE(event); err != nil {
			slog.Error(err.Error())
			return
		}
	} else {
		slog.Warn("Interaction not found", slog.String("CommandName", event.Data.CommandName()))
	}
}

func (i *interactionWorkerImpl) handleMessageCommandInteraction(event *events.ApplicationCommandInteractionCreate) {
	data := event.MessageCommandInteractionData()

	var result *MessageCommand
	for _, c := range i.manager.ListMessageCommands() {
		if data.CommandName() == c.Name {
			result = c
			break
		}
	}

	if result != nil && result.PreRunE != nil {
		if doContinue := result.PreRunE(event); !doContinue {
			return
		}
	}

	if result != nil {
		if err := result.RunE(event); err != nil {
			slog.Error(err.Error())
			return
		}
	} else {
		slog.Warn("Interaction not found", slog.String("CommandName", data.CommandName()))
	}
}

func (i *interactionWorkerImpl) handleComponentInteration(event *events.ComponentInteractionCreate) {
	var result *MessageComponent
	for _, c := range i.manager.ListMessageComponents() {
		if c.CustomIDRegex != "" {
			match, err := regexp.MatchString(c.CustomIDRegex, event.Data.CustomID())
			if match && err == nil {
				result = c
				break
			}
		}

		if c.CustomID == event.Data.CustomID() {
			result = c
			break
		}
	}

	if result != nil {
		if event.CreatedAt().Add(1*time.Second).Unix() < time.Now().Unix() {
			slog.Warn(
				fmt.Sprintf(
					"A interaction de '%s' foi recebida 1 segundo depois de ter sido criada delay? Diferença MS: %d",
					event.Data.CustomID(),
					time.Now().UnixMilli()-event.CreatedAt().UnixMilli(),
				),
			)
		}

		if err := result.RunE(event); err != nil {
			slog.Error(fmt.Sprintf("Não foi possivel completar a interaction '%s' devido ao erro: %s", event.Data.CustomID(), err.Error()))
			return
		}
	} else {
		slog.Warn("Interaction not found", slog.String("CustomID", event.Data.CustomID()))
	}
}

func (i *interactionWorkerImpl) handleModalSubmit(event *events.ModalSubmitInteractionCreate) {
	var result *ModalComponent
	for _, c := range i.manager.ListModals() {
		if c.CustomIDRegex != "" {
			match, err := regexp.MatchString(c.CustomIDRegex, event.Data.CustomID)
			if match && err == nil {
				result = c
				break
			}
		}

		if c.CustomID == event.Data.CustomID {
			result = c
			break
		}
	}

	if result != nil {
		if err := result.RunE(event); err != nil {
			slog.Error(fmt.Sprintf("Não foi possivel completar a interaction '%s' devido ao erro: %s", event.Data.CustomID, err.Error()))
			return
		}
	} else {
		slog.Warn("Interaction not found", slog.String("CustomID", event.Data.CustomID))
	}
}

func (i *interactionWorkerImpl) handleAutocompleteInteraction(event *events.AutocompleteInteractionCreate) {
	var result *SlashCommand
	for _, c := range i.manager.ListSlashCommands() {
		if event.Data.CommandName == c.Name {
			result = c
			break
		}
	}

	if result != nil && result.AutoComplete != nil {
		if err := result.AutoComplete(event); err != nil {
			slog.Error(fmt.Sprintf("Não foi possivel completar o autocomplete interaction '%s' devido ao erro: %s", event.Data.CommandName, err.Error()))
			return
		}
	} else {
		slog.Warn("Interaction not found", slog.String("CommandName", event.Data.CommandName))
	}
}

func (i *interactionWorkerImpl) Listeners() []bot.EventListener {
	return []bot.EventListener{
		bot.NewListenerFunc(func(event *events.ApplicationCommandInteractionCreate) { i.InteractionHandlerChannel <- event }),
		bot.NewListenerFunc(func(event *events.AutocompleteInteractionCreate) { i.InteractionHandlerChannel <- event }),
		bot.NewListenerFunc(func(event *events.ModalSubmitInteractionCreate) { i.InteractionHandlerChannel <- event }),
		bot.NewListenerFunc(func(event *events.ComponentInteractionCreate) { i.InteractionHandlerChannel <- event }),
	}
}
