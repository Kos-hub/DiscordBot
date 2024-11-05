package main

import (
	"discordbot/commands"
	ctx "discordbot/context"
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
	ResetFolder()
	log.Println("Closing bot...")
}

func handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	value, exists := commands.Interactions[i.ApplicationCommandData().Name]

	if !exists {
		log.Printf("Command does not exist")
	} else {
		options := i.ApplicationCommandData().Options
		var args []string
		for _, option := range options {

			args = append(args, option.StringValue())
		}
		err := value(i, args)
		if err != nil {
			log.Printf("Error with interaction '%v', '%v'", i.ApplicationCommandData().Name, err)
		}
	}

}

func ResetFolder() error {
	if err := os.RemoveAll("music"); err != nil {
		log.Printf("Error removing folder: %w", err)
		return err
	}

	if err := os.MkdirAll("music", os.ModePerm); err != nil {
		log.Printf("Error recreating folder: %w", err)
		return err
	}

	return nil
}
