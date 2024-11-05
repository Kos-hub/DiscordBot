package commands

import (
	"github.com/bwmarrin/discordgo"
)

func JoinCommand(i *discordgo.InteractionCreate) error {
	// Check if bot is in a voice channel
	if b.VoiceConnection != nil {
		b.DisplayMessage(i, "Already in voice channel!")
	}
	userID := i.Member.User.ID

	voiceState, err := b.Session.State.VoiceState(b.GuildID, userID)
	if err != nil {
		b.DisplayMessage(i, "You have to be in a voice channel for me to join.")
		return err
	}

	// If bot is not in a voice channel, then connect to the caller's voice
	vc, err := b.Session.ChannelVoiceJoin(b.GuildID, voiceState.ChannelID, false, true)
	b.VoiceConnection = vc
	if err != nil {
		b.DisplayMessage(i, "Failed to join channel")
	}

	b.DisplayMessage(i, "Joined Voice channel!")

	return nil

}
