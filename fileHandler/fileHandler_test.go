package fileHandler

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/magiconair/properties/assert"
)

func TestCreateFile(t *testing.T) {

	file := tgbotapi.FileBytes{
		Name:  "test.txt",
		Bytes: []byte("loremipsum"),
	}

	err := CreateFile(file)
	if err != nil {
		t.Fail()
	}
	_, err = os.Stat(file.Name)
	if err != nil {
		t.Fail()
	}
}

func TestDeleteFile(t *testing.T) {

	file := tgbotapi.FileBytes{
		Name:  "test.txt",
		Bytes: []byte("loremipsum"),
	}

	CreateFile(file)

	result := DeleteFile(file.Name)

	if !result {
		t.Fail()
	}
	result = DeleteFile("nothing.txt")
	assert.Equal(t, result, false)

}

func TestConvertToMp3(t *testing.T) {
	filename := "test-files/test.mpeg"
	file := tgbotapi.FileBytes{
		Name: filename,
	}
	var err error
	file.Bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	_, err = ConvertToMp3(file)
	if err != nil {
		t.Fail()
	}
	del := exec.Command("sh", "-c", "rm *.mp3").Run()
	if del != nil {
		fmt.Println(err)
		t.Fail()
	}
}
