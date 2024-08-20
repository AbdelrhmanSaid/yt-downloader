package downloader

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func DownloadVideo(videoID string) error {
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)

	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	formats := video.Formats.WithAudioChannels()

	if len(formats) == 0 {
		return fmt.Errorf("no formats available")
	}

	stream, _, err := client.GetStream(video, &formats[0])

	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}

	defer stream.Close()

	file, err := os.Create(normalizeFilename(video.Title) + ".mp4")

	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	_, err = io.Copy(file, stream)

	if err != nil {
		return fmt.Errorf("failed to copy stream to file: %w", err)
	}

	return nil
}

func DownloadAudio(videoID string) error {
	client := youtube.Client{}
	video, err := client.GetVideo(videoID)

	if err != nil {
		return fmt.Errorf("failed to get video: %w", err)
	}

	formats := video.Formats.WithAudioChannels()

	if len(formats) == 0 {
		return fmt.Errorf("no formats available")
	}

	var audioStream *youtube.Format
	for _, format := range formats {
		if strings.Contains(format.MimeType, "audio") {
			audioStream = &format
			break
		}
	}

	if audioStream == nil {
		return fmt.Errorf("no audio formats available")
	}

	stream, _, err := client.GetStream(video, audioStream)

	if err != nil {
		return fmt.Errorf("failed to get stream: %w", err)
	}

	defer stream.Close()

	file, err := os.Create(normalizeFilename(video.Title) + ".mp3")

	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	_, err = io.Copy(file, stream)

	if err != nil {
		return fmt.Errorf("failed to copy stream to file: %w", err)
	}

	return nil
}

func DownloadPlaylist(playlistID string, audio bool) error {
	client := youtube.Client{}
	playlist, err := client.GetPlaylist(playlistID)

	if err != nil {
		return fmt.Errorf("failed to get playlist: %w", err)
	}

	for _, video := range playlist.Videos {
		var err error

		if audio {
			err = DownloadAudio(video.ID)
		} else {
			err = DownloadVideo(video.ID)
		}

		if err != nil {
			return fmt.Errorf("failed to download video: %w", err)
		}
	}

	return nil
}

func normalizeFilename(filename string) string {
	return regexp.MustCompile(`[!@#$%^&*()_+=[\]{};':"\\|,.<>/?]+`).ReplaceAllString(filename, "_")
}
