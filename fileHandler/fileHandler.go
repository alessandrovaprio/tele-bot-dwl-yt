package fileHandler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ConvertToMp3(file tgbotapi.FileBytes) (string, error) {
	mp3file := strings.ReplaceAll(strings.ReplaceAll(file.Name, ".mpeg", ".mp3"), "/", "|")
	if !searchffmpeg() {
		err := errors.New("ffmpeg not found")
		return mp3file, err
	}
	cmd := exec.Command("ffmpeg", "-i", file.Name, mp3file)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	res := cmd.Run()
	if res != nil {
		fmt.Println(fmt.Sprint(res) + ": " + stderr.String())
	}
	return mp3file, res
}

// CreateFile creates a file
func CreateFile(file tgbotapi.FileBytes) error {
	localfile, errorCreateFile := os.Create(file.Name)
	if errorCreateFile != nil {
		fmt.Println(errorCreateFile)
	}
	defer localfile.Close()

	_, errorCreateFile = io.Copy(localfile, bytes.NewBuffer(file.Bytes))
	return errorCreateFile
}

func DeleteFile(filename string) bool {
	fmt.Println("delete " + filename)
	del := exec.Command("sh", "-c", "rm \""+filename+"\"").Run()
	if del != nil {
		fmt.Println(del)
		return false
	}
	return true
}

func searchffmpeg() bool {
	path, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("ffmpeg not found", path)
		return false
	}
	return true
}
