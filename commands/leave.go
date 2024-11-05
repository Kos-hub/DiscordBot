package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func LeaveCommand(i *discordgo.InteractionCreate) error {
	if b.VoiceConnection == nil {
		b.DisplayMessage(i, "I'm not in a voice channel")
		return errors.New("not in a voice channel")
	}

	b.VoiceConnection.Disconnect()

	b.DisplayMessage(i, "Left the channel!")
	b.VoiceConnection = nil // clean-up the pointer to the voice connection
	return nil
}
