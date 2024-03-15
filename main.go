package main

import (
	"fmt"
	"regexp"

	"abdelrhmansaid.com/go/youtube-downloader/downloader"
)

func main() {
	isPlaylist, _ := regexp.Compile(`list=([a-zA-Z0-9_-]+)`)

	for {
		var url string

		fmt.Print("Enter the URL of the video or playlist: ")
		fmt.Scanln(&url)

		var askForAudioOnly string

		fmt.Print("Do you want to download the audio only? (y/n): ")
		fmt.Scanln(&askForAudioOnly)

		var err error

		if isPlaylist.MatchString(url) {
			err = downloader.DownloadPlaylist(url, askForAudioOnly == "y")
		} else if askForAudioOnly == "y" {
			err = downloader.DownloadAudio(url)
		} else {
			err = downloader.DownloadVideo(url)
		}

		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Download completed successfully")
		}

		var askForAnother string

		fmt.Print("Do you want to download another video? (y/n): ")
		fmt.Scanln(&askForAnother)

		if askForAnother != "y" {
			break
		}
	}
}
