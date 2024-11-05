package commands

import (
	ctx "discordbot/context"

	"github.com/bwmarrin/discordgo"
)

var (
	Interactions map[string]func(i *discordgo.InteractionCreate, args []string) error
	b            *ctx.Bot
)

func RegisterCommands(bot *ctx.Bot) {
	Interactions = map[string]func(i *discordgo.InteractionCreate, args []string) error{
		"join":  JoinCommand,
		"leave": LeaveCommand,
		"help":  HelpCommand,
		"play":  PlayCommand,
	}
	b = bot
}
