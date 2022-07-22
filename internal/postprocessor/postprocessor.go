package postprocessor

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"

	"github.com/jawee/twitch-recorder/internal/discordclient"
	"github.com/jawee/twitch-recorder/internal/recorder"
)

type PostProcessor interface {
    Process(rf *recorder.RecordedFile) error
}

type FileMovePostProcessor struct {
    discordClient *discordclient.DiscordClient
}

func New(discordClient *discordclient.DiscordClient) *FileMovePostProcessor {
    return &FileMovePostProcessor{
        discordClient: discordClient,
    }
}


func (fm *FileMovePostProcessor) Process(rf *recorder.RecordedFile) error {
    fm.discordClient.SendMessage("Postprocessing file: " + rf.FileName)

    baseDirectory := "/videos"

    userFolderPath := path.Join(baseDirectory, rf.Username)
    if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
        os.Mkdir(userFolderPath, 0777)
    }

    filePath := path.Join(userFolderPath, rf.FileName)

    for fileExists(filePath) {
        filePath = getNewFilePath(userFolderPath, rf.FileName)
    }

    err := moveFile(rf.Path, filePath)

    if err != nil {
        log.Println("Error moving file: " + err.Error())
        return err
    }

    fm.discordClient.SendMessage(rf.FileName + " moved to processed folder")
    
    return nil
}

func getNewFilePath(userFolderPath, fileName string) string {
    re := regexp.MustCompile("^[0-9]*_[0-9]*")
    res := re.FindString(fileName)

    fileName = re.ReplaceAllString(fileName, fmt.Sprintf("%s1", res))

    return path.Join(userFolderPath, fileName)
}

func fileExists(path string) bool {
    _, error := os.Stat(path)
    return !errors.Is(error, os.ErrNotExist)
}

/*
   GoLang: os.Rename() give error "invalid cross-device link" for Docker container with Volumes.
   MoveFile(source, destination) will work moving file between folders
   https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b
*/

func moveFile(sourcePath, destPath string) error {
    inputFile, err := os.Open(sourcePath)
    if err != nil {
        return err
    }
    outputFile, err := os.Create(destPath)
    if err != nil {
        inputFile.Close()
        return err
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return err
    }
    // The copy was successful, so now delete the original file
    err = os.Remove(sourcePath)
    if err != nil {
        return err
    }
    return nil
}
