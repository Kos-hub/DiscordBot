package ctx

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Store all information I need in a struct
type Bot struct {
	Session         *discordgo.Session         // Session
	VoiceConnection *discordgo.VoiceConnection // Voice Chat
	GuildID         string                     // Guild ID <- Specific server
}

func NewBot(t string) (*Bot, error) {
	d, err := discordgo.New("Bot " + t)
	if err != nil {
		return nil, err
	}

	// Set the flags for the bot
	d.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	b := &Bot{
		Session:         d,
		VoiceConnection: nil,
		GuildID:         "523889328863051802",
	}

	return b, nil
}

func (b *Bot) DeleteCommands() {
	existingCommands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, b.GuildID)
	if err != nil {
		log.Fatalf("Did not fetch commands correctly.")
		return
	}
	for _, cmd := range existingCommands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, b.GuildID, cmd.ID)
		log.Printf("Currently deleting command: '%v', with User ID: '%v' and GuildID: '%v'", cmd.ID, b.Session.State.User.ID, b.GuildID)
		if err != nil {
			log.Fatalf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}
