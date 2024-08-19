package downloadHandler

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	fh "github.com/alessandrovaprio/tele-bot-dwl-yt/fileHandler"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kkdai/youtube/v2"
)

var (
	ytregex = regexp.MustCompile(`(http:|https:)?\/\/(www\.)?(youtube.com|youtu.be)\/(watch)?(\?v=)?(\S+)?`)
	Version = "0.0.9"
)

func DownloadMp3(url string) (tgbotapi.FileBytes, error) {
	file, err := DownloadAndConvert(url)
	fmt.Println(file.Name)
	if err != nil {
		fmt.Println(file.Name)
	}
	errFile := fh.CreateFile(file)
	if errFile == nil {
		mp3file, errConv := fh.ConvertToMp3(file)
		defer fh.DeleteFile(strings.ReplaceAll(strings.ReplaceAll(file.Name, "/", "|"), "\"", " "))
		if errConv == nil {
			file.Bytes, err = os.ReadFile(mp3file)
			file.Name = mp3file

		} else {
			err = errConv
		}
	}
	return file, err
}

// This function is going to download and REturn FileBytes
func DownloadAndConvert(url string) (tgbotapi.FileBytes, error) {

	if !ytregex.MatchString(url) {
		fmt.Println("not a youtube url")
	}
	videoID := url
	client := youtube.Client{}

	var file tgbotapi.FileBytes
	fmt.Println("before GetVideo DownloadAndConvert " + videoID)
	video, err := client.GetVideo(videoID)
	if err == nil {
		fmt.Println("before formatVideo " + videoID)
		file, err = formatVideo(video)
	}
	// if err == nil {

	// 	formats := video.Formats.WithAudioChannels()

	// 	stream, _, err := client.GetStream(video, &formats[0])

	// 	if err == nil {
	// 		tmpName := strings.ReplaceAll(video.Title+".mpeg", "/", "|")

	// 		buf := new(bytes.Buffer)
	// 		_, err = buf.ReadFrom(stream)
	// 		if err == nil {
	// 			file = tgbotapi.FileBytes{
	// 				Name:  tmpName,
	// 				Bytes: buf.Bytes(),
	// 			}
	// 		}
	// 	}
	// }
	fmt.Println("return DownloadAndConvert " + videoID)
	return file, err
}

func formatVideo(video *youtube.Video) (tgbotapi.FileBytes, error) {
	formats := video.Formats.WithAudioChannels()
	client := youtube.Client{}
	stream, _, err := client.GetStream(video, &formats[0])
	var file tgbotapi.FileBytes

	if err == nil {
		tmpName := strings.ReplaceAll(video.Title+".mpeg", "/", "|")

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(stream)
		if err == nil {
			file = tgbotapi.FileBytes{
				Name:  tmpName,
				Bytes: buf.Bytes(),
			}
		}
	}
	return file, err
}

// This function is going to download and REturn FileBytes
func DownloadPlaylist(url string) ([]tgbotapi.FileBytes, error) {

	if !ytregex.MatchString(url) {
		fmt.Println("not a youtube url")
	}
	videoID := url
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(videoID)
	files := make([]tgbotapi.FileBytes, len(playlist.Videos))
	if err == nil {

		for videoIndex := range playlist.Videos {
			entry := playlist.Videos[videoIndex]
			video, errorConv := client.VideoFromPlaylistEntry(entry)
			if errorConv == nil {
				file, error := formatVideo(video)
				if error == nil {
					files[videoIndex] = file
				}
			}
			// singlevideourl := fmt.Sprintf("%s&index=%d", url, videoIndex+1)
			// file, errorConv := DownloadAndConvert(singlevideourl)

		}
	}

	return files, err
}
