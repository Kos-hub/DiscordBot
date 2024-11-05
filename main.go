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

	defer b.Session.Close()

	log.Println("Bot is running. Press CTRL-C to exit.")

	// Wait for stop command
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	// CLEAN-UP
	b.DeleteCommands()
	log.Println("Closing bot...")
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	value, exists := commands.Interactions[i.ApplicationCommandData().Name]

	if !exists {
		log.Printf("Command does not exist")
	} else {
		err := value(i)
		if err != nil {
			log.Printf("Error with interaction '%v', '%v'", i.ApplicationCommandData().Name, err)
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
