package postprocessor

import (
	"io"
	"log"
	"os"
	"path"

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

    err := moveFile(rf.Path, filePath)

    if err != nil {
        log.Println("Error moving file: " + err.Error())
        return err
    }

    fm.discordClient.SendMessage(rf.FileName + " moved to processed folder")
    
    return nil
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
