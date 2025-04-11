package commands

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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
	sq         Queue
	sid        int32
	alive      bool
	isPlaylist bool
	Skip       bool

	mu sync.Mutex
	id int
)

func nextID() int {
	mu.Lock()
	defer mu.Unlock()
	id++
	return id
}

func PlayCommand(i *discordgo.InteractionCreate, args []string) error {

	if b.VoiceConnection == nil {
		JoinCommand(i, nil)
	}
	// Getting the link from the slash command
	if len(args) == 0 {
		return errors.New("no arguments are passed!")
	}
	b.DisplayMessage(i, args[0])

	if !alive {
		alive = true
		go playAudio()
	}

	go downloadSong(args[0])
	return nil
}

func downloadSong(url string) error {
	sid++
	log.Println("Started downlading song...")

	if strings.Contains(url, "playlist") {
		isPlaylist = true
	}

	var filename string
	if isPlaylist {
		filename = fmt.Sprintf("music/song_%(playlist_index)s.%%(ext)s", nextID())
	} else {
		filename = fmt.Sprintf("music/song_%d.%%(ext)s", nextID())
	}
	// Actual song download
	cmd := exec.Command("./yt-dlp",
		"-x",
		"--no-warning",
		"--no-progress",
		"--audio-format", "mp3",
		"--audio-quality", "0",
		"--output", filename,
		url,
	)
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command")
		return err
	}

	scanner := bufio.NewScanner(stdout)
	song := ""
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ExtractAudio") {
			idx := strings.Index(line, ":")
			if idx != -1 && idx+2 < len(line) {
				song = line[idx+2:]
				if isPlaylist {
					sq.Push(song)
				}
				continue
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error waiting... %w", err)
		return err
	}

	if !isPlaylist {
		sq.Push(song)
	}
	log.Printf("Finished download")
	return nil
}

func removeSong(song string) {
	err := os.Remove(song)

	if err != nil {
		log.Printf("Failed to delete song. ", err)
	} else {
		log.Printf("Deleted %s successfully", song)
	}
}
func playAudio() {
	for {
		done := make(chan struct{})
		if sq.IsEmpty() {
			time.Sleep(1 * time.Second)
			continue
		}
		song := sq.Pop()
		log.Printf("Song received. Now playing: %s", song)
		// Execute ffmpeg on mp3 file
		cmd := exec.Command("ffmpeg", "-i", song, "-f", "s16le", "-ar", strconv.Itoa(FRAME_RATE), "-ac", strconv.Itoa(CHANNELS), "pipe:1")

		// Read ffmpeg output
		pipe, err := cmd.StdoutPipe()

		// Create buffer for ffmpeg output
		buf := bufio.NewReaderSize(pipe, 16384)

		if err != nil {
			log.Printf("error with audio goroutine. %w", err)
			return
		}

		if err := cmd.Start(); err != nil {
			log.Printf("error with audio goroutine. %w", err)
			return
		}

		b.VoiceConnection.Speaking(true)

		// Create encoder for opus
		encoder, err := gopus.NewEncoder(FRAME_RATE, CHANNELS, gopus.Audio)

		if err != nil {
			log.Printf("Error creating encoder %w", err)
		}

		go func(done chan struct{}) {
			log.Printf("Encoder goroutine started.")
			for {
				// Create buffer for opus
				audioBuf := make([]int16, FRAME_SIZE*CHANNELS)

				// Read from ffmpeg to discord buffer
				err := binary.Read(buf, binary.LittleEndian, &audioBuf)
				if err != nil || Skip {
					log.Printf("Song is skipped or finished. Value of skip is %t and error is %w", Skip, err)
					Skip = false
					cmd.Process.Kill() // YOUNG MAN
					break
				}

				// Encode from audio buffer
				opus, err := encoder.Encode(audioBuf, FRAME_SIZE, MAX_BYTES)

				if err != nil {
					log.Fatal(err)
				}
				// Send packet back to opus
				b.VoiceConnection.OpusSend <- opus
			}

			close(done)
		}(done)

		<-done
		b.VoiceConnection.Speaking(false)
		cmd.Wait()
		removeSong(song)
	}
}
