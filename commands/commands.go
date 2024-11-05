package commands

import (
	ctx "discordbot/context"

	"github.com/bwmarrin/discordgo"
)

var (
	Interactions map[string]func(i *discordgo.InteractionCreate) error
	b            *ctx.Bot
)

func RegisterCommands(bot *ctx.Bot) {
	Interactions = map[string]func(i *discordgo.InteractionCreate) error{
		"join":  JoinCommand,
		"leave": LeaveCommand,
		"help":  HelpCommand,
	}
	b = bot
}
