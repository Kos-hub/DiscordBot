package commands

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

func LeaveCommand(i *discordgo.InteractionCreate, args []string) error {
	if b.VoiceConnection == nil {
		b.DisplayMessage(i, "Non sugnu nta nu canali vocali")
		return errors.New("not in a voice channel")
	}

	b.VoiceConnection.Disconnect()

	b.DisplayMessage(i, "Nescivi du canali")
	b.VoiceConnection = nil // clean-up the pointer to the voice connection
	return nil
}
