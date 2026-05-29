package interactivity

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

const (
	Prefix    = "panel"
	Separator = ";;"
)

type PanelManager interface {
	PanelID() string
	Render(e any) error
}

type panelManagerImpl struct {
	client *bot.Client

	ctx    context.Context
	cancel context.CancelFunc

	panelID   string
	messageID string
	channelID string
	panel     PanelModel

	resetIdle chan struct{}

	renderedAt   time.Time
	componentIds []string
}

func (m *panelManagerImpl) PanelID() string {
	return m.panelID
}

func (m *panelManagerImpl) Render(e any) error {
	v := reflect.ValueOf(e)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	field := v.FieldByName("Respond")
	if !field.IsValid() {
		return nil
	}

	responder := field.Interface().(events.InteractionResponderFunc)

	model, cmd := m.panel.Update(e)
	if cmd != nil {
		data := cmd()

		switch v := data.(type) {
		case discord.ModalCreate:
			v = v.WithCustomID(m.generateCustomID(v.CustomID))
			return responder(discord.InteractionResponseTypeModal, v)
		}
	}

	if model == nil {
		return nil
	}

	view := model.View()

	if view.Components != nil {
		for rootIdx, raw := range *view.Components {
			switch row := raw.(type) {
			case discord.SectionComponent:
				switch v := row.Accessory.(type) {
				case discord.ButtonComponent:
					fmt.Println("OK")
					row = row.WithAccessory(v.WithCustomID(m.generateCustomID(v.CustomID)))
				}

				(*view.Components)[rootIdx] = row

			case discord.ContainerComponent:
			// TODO: container component

			case discord.ActionRowComponent:
				for i, sub := range row.Components {
					switch v := sub.(type) {
					case discord.ButtonComponent:
						row.Components[i] = v.WithCustomID(m.generateCustomID(v.CustomID))
					case discord.StringSelectMenuComponent:
						row.Components[i] = v.WithCustomID(m.generateCustomID(v.CustomID))
					case discord.ChannelSelectMenuComponent:
						row.Components[i] = v.WithCustomID(m.generateCustomID(v.CustomID))
					case discord.MentionableSelectMenuComponent:
						row.Components[i] = v.WithCustomID(m.generateCustomID(v.CustomID))
					case discord.RoleSelectMenuComponent:
						row.Components[i] = v.WithCustomID(m.generateCustomID(v.CustomID))
					case discord.UserSelectMenuComponent:
						row.Components[i] = v.WithCustomID(m.generateCustomID(v.CustomID))
					}

					m.componentIds = append(m.componentIds, sub.GetCustomID())
				}

				(*view.Components)[rootIdx] = row
			}
		}
	}

	switch evt := e.(type) {
	case *events.ApplicationCommandInteractionCreate:
		msgCreate := discord.MessageCreate{
			StickerIDs:      view.StickerIDs,
			AllowedMentions: view.AllowedMentions,
			TTS:             view.TTS,
			Files:           view.Files,
		}

		if view.Content != nil {
			msgCreate.Content = *view.Content
		}

		if view.Embeds != nil {
			msgCreate.Embeds = *view.Embeds
		}

		if view.Components != nil {
			msgCreate.Components = *view.Components
		}

		if view.AllowedMentions != nil {
			msgCreate.AllowedMentions = view.AllowedMentions
		}

		if view.Flags != nil {
			msgCreate.Flags = *view.Flags
		}

		return evt.CreateMessage(msgCreate)

	// we know that the message exists through this interaction, se we can update it
	case *events.ComponentInteractionCreate, *events.ModalSubmitInteractionCreate:
		return responder(discord.InteractionResponseTypeUpdateMessage, discord.MessageUpdate{
			Content:    view.Content,
			Embeds:     view.Embeds,
			Components: view.Components,

			AllowedMentions: view.AllowedMentions,

			Flags: view.Flags,
			Files: view.Files,
		})
	}

	return nil
}

func (m *panelManagerImpl) generateCustomID(customID string) string {
	return base64.StdEncoding.EncodeToString([]byte(Prefix + Separator + m.panelID + Separator + customID))
}

func (m *panelManagerImpl) watchComponents() {
	ch, cancel := bot.NewEventCollector(m.client, func(e *events.ComponentInteractionCreate) bool {
		fmt.Println(e.Data.CustomID(), base64.RawStdEncoding.EncodeToString([]byte(Prefix+Separator+m.panelID)))
		return strings.Contains(e.Data.CustomID(), base64.RawStdEncoding.EncodeToString([]byte(Prefix+Separator+m.panelID)))
	})

	defer cancel()

	for {
		select {
		case <-m.ctx.Done():
			return
		case e := <-ch:
			go func() {
				m.Render(e)
				m.notifyActivity()
			}()
		}
	}
}

func (m *panelManagerImpl) watchModals() {
	ch, cancel := bot.NewEventCollector(m.client, func(e *events.ModalSubmitInteractionCreate) bool {
		return strings.Contains(e.Data.CustomID, base64.RawStdEncoding.EncodeToString([]byte(Prefix+Separator+m.panelID)))
	})

	defer cancel()

	for {
		select {
		case <-m.ctx.Done():
			return
		case e := <-ch:
			go func() {
				m.notifyActivity()
				m.Render(e)
			}()
		}
	}
}

func (m *panelManagerImpl) notifyActivity() {
	if m.resetIdle != nil {
		select {
		case m.resetIdle <- struct{}{}:
		default:
		}
	}
}

func (m *panelManagerImpl) monitorIdle(idleTime time.Duration) {
	timer := time.NewTimer(idleTime)
	defer timer.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-timer.C:
			m.cancel()
			return
		case <-m.resetIdle:
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

func GetCustomID(value string) string {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return ""
	}

	return strings.Split(string(decoded), Separator)[2]
}
