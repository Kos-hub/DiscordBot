package commands

import (
	"github.com/bwmarrin/discordgo"
)

func SkipCommand(i *discordgo.InteractionCreate, args []string) error {
	Skip = true

	b.DisplayMessage(i, "Skippasti a canzuni")
	return nil
}
