package main

import (
	ctx "discordbot/context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	b     ctx.Bot
)

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

	// Discord API flags.
	defer b.Session.Close()

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
				Content: "Joined Channel",
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

	if i.ApplicationCommandData().Name == "pong" {
		err := leaveVoiceChannel()

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Left Channel",
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

func leaveVoiceChannel() error {

	if b.VoiceConnection == nil {
		return errors.New("not in a voice channel")
	}

	b.VoiceConnection.Disconnect()

	b.VoiceConnection = nil // clean-up the pointer to the voice connection
	return nil
}

func joinUserVoiceChannel(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	gID := i.GuildID
	userID := i.Member.User.ID

	voiceState, err := s.State.VoiceState(gID, userID)
	if err != nil {
		return fmt.Errorf("could not retrieve voice: %v", err)
	}

	if voiceState == nil || voiceState.ChannelID == "" {
		return fmt.Errorf("you need to be in a voice channel for me to join")
	}

	vc, err := s.ChannelVoiceJoin(gID, voiceState.ChannelID, false, true)
	b.VoiceConnection = vc
	if err != nil {
		return fmt.Errorf("failed to join voice channel: %v", err)
	}

	return nil
}
