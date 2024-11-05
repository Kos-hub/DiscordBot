package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func HelpCommand(i *discordgo.InteractionCreate) error {

	builder := new(strings.Builder)
	for _, c := range b.Commands {
		builder.WriteString(fmt.Sprintf("/%s - %s\n", c.Name, c.Description))
	}

	if builder.Len() == 0 {
		b.DisplayMessage(i, "No commands have been found")
		return errors.New("could not find any commands to display")
	}
	b.DisplayMessage(i, builder.String())
	return nil
}
