package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/richaardev/concord/interactivity"
)

type User struct {
	ID   int
	Name string
	Age  int
}
type UsersView struct {
	users       []User
	nextID      int
	confirmID   int
	editUserIdx int
}

func (s *UsersView) Init() {}
func (s *UsersView) Update(e any) (interactivity.PanelModel, interactivity.Cmd) {
	switch i := e.(type) {
	case *events.ComponentInteractionCreate:
		return s.handleComponent(i)
	case *events.ModalSubmitInteractionCreate:
		return s.handleModal(i)
	}
	return s, nil
}

func (s *UsersView) handleComponent(e *events.ComponentInteractionCreate) (interactivity.PanelModel, interactivity.Cmd) {
	customID := interactivity.GetCustomID(e.Data.CustomID())
	switch {
	case customID == "users:add":
		return nil, interactivity.ModalCmd(discord.ModalCreate{
			Title:    "Adicionar usuário",
			CustomID: "modal:user:add",
			Components: []discord.LayoutComponent{
				discord.NewLabel(
					"Nome",
					discord.NewShortTextInput("user_name").WithRequired(true).WithPlaceholder("Nome do usuário"),
				),
				discord.NewLabel(
					"Idade",
					discord.NewShortTextInput("user_age").WithRequired(true).WithPlaceholder("Idade"),
				),
			},
		})
	case customID == "users:confirm_remove":
		idx := -1
		for i, u := range s.users {
			if u.ID == s.confirmID {
				idx = i
				break
			}
		}
		if idx >= 0 {
			s.users = append(s.users[:idx], s.users[idx+1:]...)
		}
		s.confirmID = 0
		return s, nil
	case customID == "users:cancel_remove":
		s.confirmID = 0
		return s, nil
	case strings.HasPrefix(customID, "users:remove:"):
		idStr := strings.TrimPrefix(customID, "users:remove:")
		id, _ := strconv.Atoi(idStr)
		s.confirmID = id
		return s, nil
	case strings.HasPrefix(customID, "users:edit:"):
		idStr := strings.TrimPrefix(customID, "users:edit:")
		id, _ := strconv.Atoi(idStr)
		for i, u := range s.users {
			if u.ID == id {
				s.editUserIdx = i
				return nil, interactivity.ModalCmd(discord.ModalCreate{
					Title:    fmt.Sprintf("Editar %s", u.Name),
					CustomID: "modal:user:edit",
					Components: []discord.LayoutComponent{
						discord.NewLabel(
							"Nome",
							discord.NewShortTextInput("user_name").WithRequired(true).WithValue(u.Name),
						),
						discord.NewLabel(
							"Idade",
							discord.NewShortTextInput("user_age").WithRequired(true).WithValue(strconv.Itoa(u.Age)),
						),
					},
				})
			}
		}
	}
	return s, nil
}

func (s *UsersView) handleModal(e *events.ModalSubmitInteractionCreate) (interactivity.PanelModel, interactivity.Cmd) {
	customID := interactivity.GetCustomID(e.Data.CustomID)
	name := e.Data.Text("user_name")
	ageStr := e.Data.Text("user_age")
	age, _ := strconv.Atoi(ageStr)
	switch customID {
	case "modal:user:add":
		s.users = append(s.users, User{ID: s.nextID, Name: name, Age: age})
		s.nextID++
	case "modal:user:edit":
		if s.editUserIdx >= 0 && s.editUserIdx < len(s.users) {
			s.users[s.editUserIdx].Name = name
			s.users[s.editUserIdx].Age = age
		}
	}
	return s, nil
}

func (s *UsersView) View() interactivity.Message {
	if s.confirmID > 0 {
		return s.confirmRemoveView()
	}
	return s.listView()
}

func (s *UsersView) listView() interactivity.Message {
	var sections []discord.LayoutComponent
	if len(s.users) == 0 {
		sections = append(
			sections,
			discord.NewSection(discord.NewTextDisplay("Nenhum usuário cadastrado.")),
		)
	}
	for _, u := range s.users {
		idStr := strconv.Itoa(u.ID)
		sections = append(
			sections,
			discord.NewSection(
				discord.NewTextDisplay(
					strings.Join([]string{
						fmt.Sprintf("**%s** (ID: %d)", u.Name, u.ID),
						fmt.Sprintf("Idade: %d", u.Age),
					}, "\n"),
				),
			).WithAccessory(
				discord.NewSecondaryButton("✏️", fmt.Sprintf("users:edit:%s", idStr)),
			).WithAccessory(
				discord.NewDangerButton("🗑️", fmt.Sprintf("users:remove:%s", idStr)),
			),
		)
	}

	sections = append(
		sections,
		discord.NewActionRow(
			discord.NewSuccessButton("+ Adicionar", "users:add"),
		),
	)

	return interactivity.Message{
		Flags:      new(discord.MessageFlagEphemeral | discord.MessageFlagIsComponentsV2),
		Components: new(sections),
	}
}

func (s *UsersView) confirmRemoveView() interactivity.Message {
	return interactivity.Message{
		Flags: new(discord.MessageFlagEphemeral | discord.MessageFlagIsComponentsV2),
		Components: &[]discord.LayoutComponent{
			discord.NewSection(
				discord.NewTextDisplay("Tem certeza que deseja remover este usuário?"),
			),
			discord.NewActionRow(
				discord.NewDangerButton("Sim, remover", "users:confirm_remove"),
				discord.NewSecondaryButton("Cancelar", "users:cancel_remove"),
			),
		},
	}
}
