package concord

var _ InteractionManager = (*interactionManagerImpl)(nil)

type InteractionManager interface {
	ListSlashCommands() []*SlashCommand
	GetSlashCommand(name string) *SlashCommand
	RegisterSlashCommands(interactives ...*SlashCommand)

	ListModals() []*ModalComponent
	GetModal(id string) *ModalComponent
	RegisterModals(interactives ...*ModalComponent)

	ListMessageCommands() []*MessageCommand
	GetMessageCommand(name string) *MessageCommand
	RegisterMessageCommands(interactives ...*MessageCommand)

	ListMessageComponents() []*MessageComponent
	GetMessageComponents(customId string) []*MessageComponent
	RegisterMessageComponents(interactives ...*MessageComponent)
}

type interactionManagerImpl struct {
	modals            []*ModalComponent
	slashCommands     []*SlashCommand
	messageCommands   []*MessageCommand
	messageComponents []*MessageComponent
}

func NewInteractionManager() InteractionManager {
	return &interactionManagerImpl{}
}

func (manager *interactionManagerImpl) RegisterModals(interactives ...*ModalComponent) {
	manager.modals = append(manager.modals, interactives...)
}

func (manager *interactionManagerImpl) RegisterSlashCommands(interactives ...*SlashCommand) {
	manager.slashCommands = append(manager.slashCommands, interactives...)
}

func (manager *interactionManagerImpl) RegisterMessageCommands(interactives ...*MessageCommand) {
	manager.messageCommands = append(manager.messageCommands, interactives...)
}

func (manager *interactionManagerImpl) RegisterMessageComponents(interactives ...*MessageComponent) {
	manager.messageComponents = append(manager.messageComponents, interactives...)
}

func (manager *interactionManagerImpl) ListSlashCommands() []*SlashCommand {
	return manager.slashCommands
}

func (manager *interactionManagerImpl) GetSlashCommand(name string) *SlashCommand {
	for _, cmd := range manager.slashCommands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

func (manager *interactionManagerImpl) ListModals() []*ModalComponent {
	return manager.modals
}

func (manager *interactionManagerImpl) GetModal(id string) *ModalComponent {
	for _, modal := range manager.modals {
		if modal.CustomID == id {
			return modal
		}
	}
	return nil
}

func (manager *interactionManagerImpl) ListMessageCommands() []*MessageCommand {
	return manager.messageCommands
}

func (manager *interactionManagerImpl) GetMessageCommand(name string) *MessageCommand {
	for _, cmd := range manager.messageCommands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

func (manager *interactionManagerImpl) ListMessageComponents() []*MessageComponent {
	return manager.messageComponents
}

func (manager *interactionManagerImpl) GetMessageComponents(customId string) []*MessageComponent {
	var matches []*MessageComponent
	for _, comp := range manager.messageComponents {
		if comp.CustomID == customId {
			matches = append(matches, comp)
		}
	}
	return matches
}
