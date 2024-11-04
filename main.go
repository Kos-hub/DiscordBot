package main

import (
	ctx "discordbot/context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var Token string

// Arguments that are passed in when calling go-run
func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Initialize the bot
	b, err := ctx.NewBot(Token)
	if err != nil {
		log.Fatalln("Error creating discord session,", err)
		return
	}

	// Event handlers.

	b.Session.AddHandler(handleInteraction)

	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "ping",
			Description: "Ping!",
		},
	}
	// Discord API flags.
	err = b.Session.Open()
	if err != nil {
		log.Fatalln("Error opening connection,", err)
	}
	defer b.Session.Close()

	b.DeleteCommands()

	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, b.GuildID, cmd)
		log.Printf("Currently adding command: '%v', with User ID: '%v' and GuildID: '%v'", cmd.ID, b.Session.State.User.ID, b.GuildID)
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", cmd.Name, err)
		}
	}
	log.Println("Bot is running. Press CTRL-C to exit.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	b.DeleteCommands()
	log.Println("Closing bot...")
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.ApplicationCommandData().Name == "ping" {
		err := joinUserVoiceChannel(s, i)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Joining voice channel...",
			},
		})

		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Error: " + err.Error(),
				},
			})

			return
		}
	}

}

func joinUserVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	gID := i.GuildID
	userID := i.Member.User.ID

	voiceState, err := s.State.VoiceState(gID, userID)
	if err != nil {
		return fmt.Errorf("could not retrieve voice: %v", err)
	}

	if voiceState == nil || voiceState.ChannelID == "" {
		return fmt.Errorf("you need to be in a voice channel for me to join!")
	}

	_, err = s.ChannelVoiceJoin(gID, voiceState.ChannelID, false, true)
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %v", err)
	}

	return nil
}