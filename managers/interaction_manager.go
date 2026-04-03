package managers

import "github.com/richaardev/concord/discord/core"

var _ InteractionManager = (*interactionManagerImpl)(nil)

type InteractionManager interface {
	RegisterModals(interactives ...*core.ModalComponent)
	RegisterSlashCommands(interactives ...*core.SlashCommand)
	RegisterMessageCommands(interactives ...*core.MessageCommand)
	RegisterMessageComponents(interactives ...*core.MessageComponent)

	Modals() []*core.ModalComponent
	SlashCommands() []*core.SlashCommand
	MessageCommands() []*core.MessageCommand
	MessageComponents() []*core.MessageComponent
}

type interactionManagerImpl struct {
	modals            []*core.ModalComponent
	slashCommands     []*core.SlashCommand
	messageCommands   []*core.MessageCommand
	messageComponents []*core.MessageComponent
}

func NewInteractionManager() InteractionManager {
	return &interactionManagerImpl{}
}

func (manager *interactionManagerImpl) RegisterModals(interactives ...*core.ModalComponent) {
	manager.modals = append(manager.modals, interactives...)
}
func (manager *interactionManagerImpl) RegisterSlashCommands(interactives ...*core.SlashCommand) {
	manager.slashCommands = append(manager.slashCommands, interactives...)
}
func (manager *interactionManagerImpl) RegisterMessageCommands(interactives ...*core.MessageCommand) {
	manager.messageCommands = append(manager.messageCommands, interactives...)
}
func (manager *interactionManagerImpl) RegisterMessageComponents(interactives ...*core.MessageComponent) {
	manager.messageComponents = append(manager.messageComponents, interactives...)
}

func (manager *interactionManagerImpl) Modals() []*core.ModalComponent {
	return manager.modals
}

func (manager *interactionManagerImpl) SlashCommands() []*core.SlashCommand {
	return manager.slashCommands
}

func (manager *interactionManagerImpl) MessageCommands() []*core.MessageCommand {
	return manager.messageCommands
}

func (manager *interactionManagerImpl) MessageComponents() []*core.MessageComponent {
	return manager.messageComponents
}
