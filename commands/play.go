package commands

import (
	"errors"
	"log"
	"os/exec"

	"github.com/bwmarrin/discordgo"
)

func PlayCommand(i *discordgo.InteractionCreate, args []string) error {

	// Getting the link from the slash command
	if len(args) == 0 {
		return errors.New("no arguments are passed!")
	}
	b.DisplayMessage(i, args[0])

	// Download the song from ytdl
	downloadSong(args[0])

	// Reproduce music

	return nil
}

func downloadSong(url string) error {
	cmd := exec.Command("./yt-dlp",
		"-x",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--embed-thumbnail",
		"--output", "music/%(title)s.%(ext)s",
		url,
	)

	err := cmd.Run()

	if err != nil {
		log.Printf("Error downloading song: %v", err)
		return err
	}

	return nil
}
