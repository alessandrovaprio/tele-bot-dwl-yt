package downloadHandler

import (
	"bytes"
	"fmt"
	"io/ioutil"
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
		fh.DeleteFile(strings.ReplaceAll(file.Name, "/", "|"))
		if errConv == nil {
			file.Bytes, err = ioutil.ReadFile(mp3file)
			file.Name = mp3file
			// in any case I delete the mp3 file
			fh.DeleteFile(mp3file)
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

	video, err := client.GetVideo(videoID)
	if err == nil {

		formats := video.Formats.WithAudioChannels()

		stream, _, err := client.GetStream(video, &formats[0])

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
	}

	return file, err
}
