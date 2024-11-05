package main

import (
	"discordbot/commands"
	ctx "discordbot/context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
	b     *ctx.Bot
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

	commands.RegisterCommands(b)

	// Event handlers.
	b.Session.AddHandler(handleInteraction)

	// Discord API flags.
	defer b.Session.Close()

	log.Println("Bot is running. Press CTRL-C to exit.")

	if b != nil {
		log.Println("Bot is not null at this point.")
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	b.DeleteCommands()
	log.Println("Closing bot...")
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//
	//	if i.ApplicationCommandData().Name == "ping" {
	//		err := joinUserVoiceChannel(s, i)
	//
	//		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: "Joined Channel",
	//			},
	//		})
	//
	//		if err != nil {
	//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//				Type: discordgo.InteractionResponseChannelMessageWithSource,
	//				Data: &discordgo.InteractionResponseData{
	//					Content: "Error: " + err.Error(),
	//				},
	//			})
	//
	//			return
	//		}
	//	}
	//
	//	if i.ApplicationCommandData().Name == "pong" {
	//		err := leaveVoiceChannel()
	//
	//		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//			Type: discordgo.InteractionResponseChannelMessageWithSource,
	//			Data: &discordgo.InteractionResponseData{
	//				Content: "Left Channel",
	//			},
	//		})
	//
	//		if err != nil {
	//			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	//				Type: discordgo.InteractionResponseChannelMessageWithSource,
	//				Data: &discordgo.InteractionResponseData{
	//					Content: "Error: " + err.Error(),
	//				},
	//			})
	//
	//			return
	//		}
	//	}

	value, exists := commands.Interactions[i.ApplicationCommandData().Name]

	if !exists {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Command doesn't exist",
			},
		})
	} else {
		value(i)
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
