package commands

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	CHANNELS   int = 2
	FRAME_RATE int = 48000
	FRAME_SIZE int = 960
	MAX_BYTES  int = (FRAME_SIZE * 2) * 2
)

var buffer = make([][]byte, 0)

func PlayCommand(i *discordgo.InteractionCreate, args []string) error {

	// Getting the link from the slash command
	if len(args) == 0 {
		return errors.New("no arguments are passed!")
	}
	b.DisplayMessage(i, args[0])

	// Download the song from ytdl
	downloadSong(args[0])

	// Reproduce music

	go playAudio()

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

func getFirstFile() string {

	files, err := os.ReadDir("music")
	if err != nil {
		log.Printf("Error opening folder, %w", err)
		return ""
	}

	for _, file := range files {
		if !file.IsDir() {
			return filepath.Join("music", file.Name())
		}
	}

	return ""
}

func playAudio() {
	cmd := exec.Command("ffmpeg", "-i", getFirstFile(), "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")

	pipe, err := cmd.StdoutPipe()
	buf := bufio.NewReaderSize(pipe, 16384)
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	b.VoiceConnection.Speaking(true)
	go func() {
		encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
		if err != nil {
			log.Printf("Error creating encoder %w", err)
		}
		for {

			audioBuf := make([]int16, FRAME_SIZE*CHANNELS)
			err := binary.Read(buf, binary.LittleEndian, &audioBuf)
			if err != nil {
				break
			}

			opus, err := encoder.Encode(audioBuf, FRAME_SIZE, MAX_BYTES)

			b.VoiceConnection.OpusSend <- opus
		}
	}()
	b.VoiceConnection.Speaking(false)

	cmd.Wait()
}
