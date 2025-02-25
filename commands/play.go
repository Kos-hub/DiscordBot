package commands

import (
	"bufio"
	"encoding/binary"
	"errors"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"layeh.com/gopus"
)

const (
	CHANNELS   int = 2
	FRAME_RATE int = 48000
	FRAME_SIZE int = 960
	MAX_BYTES  int = (FRAME_SIZE * 2) * 2
)

var (
	dq    Queue
	alive bool
)

func PlayCommand(i *discordgo.InteractionCreate, args []string) error {
	// Getting the link from the slash command
	if len(args) == 0 {
		return errors.New("no arguments are passed!")
	}
	b.DisplayMessage(i, args[0])

	if !alive {
		dq.Push(args[0])
		playChan := make(chan string, 10)
		alive = true
		go downloadQueue(playChan)
		go playAudio(playChan)
	} else {
		dq.Push(args[0])
	}

	return nil
}

func downloadQueue(playChan chan string) {
	for {
		if dq.IsEmpty() {
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Queue is not empty. about to download song...")
		title, err := downloadSong(dq.Pop())
		if err != nil {
			continue
		}

		playChan <- title
		<-playChan

	}
}

func downloadSong(url string) (string, error) {
	log.Println("Started downlading song...")

	// Actual song download
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
		return "", err
	}

	// Second call to retreive title - Might be able to get metadata from prev command
	output, err := exec.Command("./yt-dlp",
		"--print", "title",
		"--no-warning",
		url,
	).CombinedOutput()

	if err != nil {
		log.Printf("Error printing title: %v", err)
	}

	log.Printf("------------")
	return string(output[:]), nil
}

func playAudio(playChan chan string) {
	log.Printf("Started audio goroutine.")
	for s := range playChan {
		log.Printf("Song received.")
		// Execute ffmpeg on mp3 file
		song := strings.TrimSuffix(s, "\n")
		cmd := exec.Command("ffmpeg", "-i", "music\\"+song+".mp3", "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")

		// Read ffmpeg output
		pipe, err := cmd.StdoutPipe()

		// Create buffer for ffmpeg output
		buf := bufio.NewReaderSize(pipe, 16384)

		if err != nil {
			return
		}

		if err := cmd.Start(); err != nil {
			return
		}

		b.VoiceConnection.Speaking(true)

		// Create encoder for opus
		encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)

		if err != nil {
			log.Printf("Error creating encoder %w", err)
		}

		go func() {
			for {
				// Create buffer for opus
				audioBuf := make([]int16, FRAME_SIZE*CHANNELS)

				// Read from ffmpeg to discord buffer
				err := binary.Read(buf, binary.LittleEndian, &audioBuf)
				if err != nil {
					playChan <- ""
					break
				}

				// Encode from audio buffer
				opus, err := encoder.Encode(audioBuf, FRAME_SIZE, MAX_BYTES)

				// Send packet back to opus
				b.VoiceConnection.OpusSend <- opus
			}
		}()
		b.VoiceConnection.Speaking(false)

		cmd.Wait()
	}
	log.Printf("Finished goroutine.")
}
