package ctx

import (
	"encoding/json"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

// Store all information I need in a struct
type Bot struct {
	Session         *discordgo.Session         // Session
	VoiceConnection *discordgo.VoiceConnection // Voice Chat
	GuildID         string                     // Guild ID <- Specific server
	commands        []*discordgo.ApplicationCommand
}

func NewBot(t string) (*Bot, error) {
	// Discord API to retreive the session
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

	err = b.Session.Open()
	if err != nil {
		log.Fatalln("Error opening connection,", err)
	}

	b.addCommandsJSON()
	b.addSlashCommands()

	return b, nil
}

func (b *Bot) addCommandsJSON() {
	data, err := os.ReadFile("commands.json")
	if err != nil {
		log.Fatalln("Could not load JSON,", err)
	}

	err = json.Unmarshal(data, &b.commands)
	if err != nil {
		log.Fatalln("Could not unmarshal JSON,", err)
	}
}

func (b *Bot) addSlashCommands() {
	// This will send a POST request to register commands
	for _, cmd := range b.commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
		log.Printf("Currently adding command: '%v', with User ID: '%v' and GuildID: '%v'", cmd.ID, b.Session.State.User.ID, b.GuildID)
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", cmd.Name, err)
		}
	}
}

func (b *Bot) DeleteCommands() {
	// Fetch all commands and delete them
	existingCommands, err := b.Session.ApplicationCommands(b.Session.State.User.ID, b.GuildID)
	if err != nil {
		log.Fatalf("Did not fetch commands correctly.", err)
	}
	for _, cmd := range existingCommands {
		err := b.Session.ApplicationCommandDelete(b.Session.State.User.ID, b.GuildID, cmd.ID)
		log.Printf("Currently deleting command: '%v', with User ID: '%v' and GuildID: '%v'", cmd.ID, b.Session.State.User.ID, b.GuildID)
		if err != nil {
			log.Fatalf("Cannot delete '%v' command: %v", cmd.Name, err)
		}
	}
}

// Simple abstraction to easily print messages on discord.
func (b *Bot) DisplayMessage(i *discordgo.InteractionCreate, c string) {
	b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: c,
		},
	})
}
