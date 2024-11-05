package commands

import (
	"github.com/bwmarrin/discordgo"
)

func JoinCommand(i *discordgo.InteractionCreate, args []string) error {
	// Check if bot is in a voice channel
	if b.VoiceConnection != nil {
		b.DisplayMessage(i, "Sugnu gia' nto canali")
	}
	userID := i.Member.User.ID

	voiceState, err := b.Session.State.VoiceState(b.GuildID, userID)
	if err != nil {
		b.DisplayMessage(i, "Ndai u si nta nu canali vocali pemmu pozzu trasiri")
		return err
	}

	// If bot is not in a voice channel, then connect to the caller's voice
	vc, err := b.Session.ChannelVoiceJoin(b.GuildID, voiceState.ChannelID, false, true)
	b.VoiceConnection = vc
	if err != nil {
		b.DisplayMessage(i, "Non pozzu trasiri")
	}

	b.DisplayMessage(i, "Trasivi nto canali vocali")

	return nil

}
