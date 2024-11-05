package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func JoinCommand(i *discordgo.InteractionCreate) error {
	// Check if bot is in a voice channel
	//if b.VoiceConnection != nil {
	//	log.Fatalf("Already in a channel!")
	//}
	userID := i.Member.User.ID

	voiceState, err := b.Session.State.VoiceState(b.GuildID, userID)
	if err != nil {
		log.Fatalf("could not retrieve voice: %v", err)
	}

	if voiceState == nil || voiceState.ChannelID == "" {
		log.Fatalf("you need to be in a voice channel for me to join")
	}

	// If bot is not in a voice channel, then connect to the caller's voice
	vc, err := b.Session.ChannelVoiceJoin(b.GuildID, voiceState.ChannelID, false, true)
	b.VoiceConnection = vc
	if err != nil {
		log.Fatalf("failed to join voice channel: %v", err)
	}
	// Send reply to confirm

	return nil

}
