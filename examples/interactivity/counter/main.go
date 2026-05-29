package main

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/richaardev/concord/interactivity"
)

type CounterView struct {
	count int
}

func (s *CounterView) Init() {}
func (s *CounterView) Update(e any) (interactivity.PanelModel, interactivity.Cmd) {
	switch i := e.(type) {
	case *events.ComponentInteractionCreate:
		customID := interactivity.GetCustomID(i.Data.CustomID())
		switch customID {
		case "counter:inc":
			s.count++
			return s, nil
		case "counter:dec":
			s.count--
			return s, nil
		case "counter:reset":
			s.count = 0
			return s, nil
		}
	}
	return s, nil
}

func (s *CounterView) View() interactivity.Message {
	return interactivity.Message{
		Flags: new(discord.MessageFlagEphemeral | discord.MessageFlagIsComponentsV2),
		Components: &[]discord.LayoutComponent{
			discord.NewActionRow(
				discord.NewSecondaryButton("-", "counter:dec"),
				discord.NewSecondaryButton(fmt.Sprintf("%d", s.count), "counter:nop").WithDisabled(true),
				discord.NewSecondaryButton("+", "counter:inc"),
				discord.NewDangerButton("Resetar", "counter:reset"),
			),
		},
	}
}
