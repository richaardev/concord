package discordutils

import (
	"fmt"
	"slices"

	"github.com/disgoorg/disgo/discord"
)

func UserMention(id string) string {
	return fmt.Sprintf("<@!%s>", id)
}

func RoleMention(id string) string {
	return fmt.Sprintf("<@&%s>", id)
}

func ChannelMention(id string) string {
	return fmt.Sprintf("<#%s>", id)
}

func ToMessageCreate(msg discord.Message) discord.MessageCreate {
	return discord.MessageCreate{
		Content:    msg.Content,
		TTS:        msg.TTS,
		Embeds:     slices.Clone(msg.Embeds),
		Components: slices.Clone(msg.Components),

		AllowedMentions: nil,
		Nonce:           string(msg.Nonce),
		Flags:           msg.Flags,
	}
}

func ToMessageUpdate(create discord.MessageCreate) discord.MessageUpdate {
	var update discord.MessageUpdate

	if create.Content != "" {
		content := create.Content
		update.Content = &content
	}

	if len(create.Embeds) > 0 {
		embeds := slices.Clone(create.Embeds)
		update.Embeds = &embeds
	}

	if len(create.Components) > 0 {
		components := slices.Clone(create.Components)
		update.Components = &components
	}

	if len(create.Attachments) > 0 {
		attachments := make([]discord.AttachmentUpdate, len(create.Attachments))
		for i, attachment := range create.Attachments {
			attachments[i] = attachment
		}
		update.Attachments = &attachments
	}

	if len(create.Files) > 0 {
		update.Files = append(make([]*discord.File, 0, len(create.Files)), create.Files...)
	}

	if create.AllowedMentions != nil {
		update.AllowedMentions = create.AllowedMentions
	}

	if create.Flags != discord.MessageFlagsNone {
		flags := create.Flags
		update.Flags = &flags
	}

	return update
}

func DisableMessage(message discord.MessageCreate) discord.MessageCreate {
	message.Components = DisableMessageComponents(message.Components)
	return message
}

func DisableMessageComponents(components []discord.LayoutComponent) []discord.LayoutComponent {
	result := make([]discord.LayoutComponent, len(components))

	for i, component := range components {
		result[i] = disableLayoutComponent(component)
	}

	return result
}

func disableLayoutComponent(component discord.LayoutComponent) discord.LayoutComponent {
	switch c := component.(type) {
	case discord.ActionRowComponent:
		return disableActionRowComponents(c)
	case discord.SectionComponent:
		return disableSectionComponent(c)
	case discord.ContainerComponent:
		return disableContainerComponent(c)
	default:
		return component
	}
}

func disableActionRowComponents(row discord.ActionRowComponent) discord.ActionRowComponent {
	disabledComponents := make([]discord.InteractiveComponent, len(row.Components))

	for i, component := range row.Components {
		disabledComponents[i] = disableInteractiveComponent(component)
	}

	return row.WithComponents(disabledComponents...)
}

func disableInteractiveComponent(component discord.InteractiveComponent) discord.InteractiveComponent {
	switch c := component.(type) {
	case discord.ButtonComponent:
		return c.AsDisabled()
	case discord.StringSelectMenuComponent:
		return c.AsDisabled()
	case discord.UserSelectMenuComponent:
		return c.AsDisabled()
	case discord.RoleSelectMenuComponent:
		return c.AsDisabled()
	case discord.MentionableSelectMenuComponent:
		return c.AsDisabled()
	case discord.ChannelSelectMenuComponent:
		return c.AsDisabled()
	default:
		return component
	}
}

func disableSectionComponent(section discord.SectionComponent) discord.SectionComponent {
	if section.Accessory == nil {
		return section
	}

	switch c := section.Accessory.(type) {
	case discord.ButtonComponent:
		return section.WithAccessory(c.AsDisabled())
	default:
		return section
	}
}

func disableContainerComponent(container discord.ContainerComponent) discord.ContainerComponent {
	disabledComponents := make([]discord.ContainerSubComponent, len(container.Components))

	for i, component := range container.Components {
		disabledComponents[i] = disableContainerSubComponent(component)
	}

	return container.WithComponents(disabledComponents...)
}

func disableContainerSubComponent(component discord.ContainerSubComponent) discord.ContainerSubComponent {
	switch c := component.(type) {
	case discord.ActionRowComponent:
		return disableActionRowComponents(c)
	case discord.SectionComponent:
		return disableSectionComponent(c)
	default:
		return component
	}
}

func ValidateDiscordMessage(msg discord.MessageCreate) error {
	if len(msg.Content) > 2000 {
		return fmt.Errorf("conteúdo excede 2000 caracteres (%d caracteres)", len(msg.Content))
	}

	if len(msg.Embeds) > 10 {
		return fmt.Errorf("máximo de 10 embeds permitidos (%d fornecidos)", len(msg.Embeds))
	}

	for i, embed := range msg.Embeds {
		totalChars := 0

		if len(embed.Title) > 256 {
			return fmt.Errorf("embed %d: título excede 256 caracteres (%d)", i+1, len(embed.Title))
		}
		totalChars += len(embed.Title)

		if len(embed.Description) > 4096 {
			return fmt.Errorf("embed %d: descrição excede 4096 caracteres (%d)", i+1, len(embed.Description))
		}
		totalChars += len(embed.Description)

		if len(embed.Fields) > 25 {
			return fmt.Errorf("embed %d: máximo de 25 fields permitidos (%d fornecidos)", i+1, len(embed.Fields))
		}

		for j, field := range embed.Fields {

			if len(field.Name) > 256 {
				return fmt.Errorf("embed %d, field %d: nome excede 256 caracteres (%d)", i+1, j+1, len(field.Name))
			}
			totalChars += len(field.Name)

			if len(field.Value) > 1024 {
				return fmt.Errorf("embed %d, field %d: valor excede 1024 caracteres (%d)", i+1, j+1, len(field.Value))
			}
			totalChars += len(field.Value)
		}

		if embed.Footer != nil && len(embed.Footer.Text) > 2048 {
			return fmt.Errorf("embed %d: footer excede 2048 caracteres (%d)", i+1, len(embed.Footer.Text))
		}
		if embed.Footer != nil {
			totalChars += len(embed.Footer.Text)
		}

		if embed.Author != nil && len(embed.Author.Name) > 256 {
			return fmt.Errorf("embed %d: nome do author excede 256 caracteres (%d)", i+1, len(embed.Author.Name))
		}
		if embed.Author != nil {
			totalChars += len(embed.Author.Name)
		}

		if totalChars > 6000 {
			return fmt.Errorf("embed %d: total de caracteres excede 6000 (%d)", i+1, totalChars)
		}
	}

	return nil
}
