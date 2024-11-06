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

var queue Queue

func PlayCommand(i *discordgo.InteractionCreate, args []string) error {
	// Getting the link from the slash command
	if len(args) == 0 {
		return errors.New("no arguments are passed!")
	}
	b.DisplayMessage(i, args[0])
	queue.Push(args[0])

	for !queue.IsEmpty() {
		// Download the song from ytdl
		log.Printf("Queue is not empty. about to download song...")
		downloadSong(queue.Pop())
	}

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

	log.Println("Started downlading song...")
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
	// Execute ffmpeg on mp3 file
	cmd := exec.Command("ffmpeg", "-i", getFirstFile(), "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")

	log.Println("Starting ffmpeg data streaming.")

	// Read ffmpeg output
	pipe, err := cmd.StdoutPipe()

	// Create buffer for ffmpeg output
	buf := bufio.NewReaderSize(pipe, 16384)

	log.Println("Created buffer for streaming")
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	b.VoiceConnection.Speaking(true)

	// Create encoder for opus
	encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)
	log.Println("Created encoder")

	if err != nil {
		log.Printf("Error creating encoder %w", err)
	}

	go func() {
		for {
			// Create buffer for opus
			audioBuf := make([]int16, FRAME_SIZE*CHANNELS)

			// Read from ffmpeg to discord buffer
			err := binary.Read(buf, binary.LittleEndian, &audioBuf)
			log.Println("Read from buffer. Audio buffer is now:", audioBuf)
			if err != nil {
				break
			}

			// Encode from audio buffer
			opus, err := encoder.Encode(audioBuf, FRAME_SIZE, MAX_BYTES)

			log.Println("Encoded OPUS string is:", opus)
			// Send packet back to opus
			b.VoiceConnection.OpusSend <- opus
		}
	}()
	b.VoiceConnection.Speaking(false)

	cmd.Wait()
}
