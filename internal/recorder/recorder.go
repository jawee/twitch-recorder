package recorder

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/jawee/twitch-recorder/internal/discordclient"
)

type RecordedFile struct {
    Username string
    FileName string
    Path string
}

type Recorder interface {
    Record(username string, filename string) (*RecordedFile, error)
}

type StreamlinkRecorder struct {
    baseDirectory string
    notificationClient discordclient.NotificationClient
}

func New(baseDirectory string, notificationClient discordclient.NotificationClient) *StreamlinkRecorder {
    return &StreamlinkRecorder{
        baseDirectory: baseDirectory,
        notificationClient: notificationClient,
    }
}

func (s* StreamlinkRecorder) Record(username string, filename string) (*RecordedFile, error) {
    log.Println("Starting recording")

    filePath := path.Join(s.baseDirectory, username, filename)

    _, err := os.Stat(filePath)
    if err == nil {
        return nil, errors.New("File already exists")
    }

    s.notificationClient.SendMessage("Starting recording for " + username + ". File " + filename)

    userFolderPath := path.Join(s.baseDirectory, username)
    if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
        os.Mkdir(userFolderPath, 0777)
    }
    cmd := exec.Command("streamlink", "twitch.tv/" + username, "best", "-o", filePath)

    log.Println("Running cmd")
    cmd.Stdout = os.Stdout
    err = cmd.Run()

    if err != nil {
        log.Println(err)
        return nil, err
    }


    return &RecordedFile{
        Username: username,
        FileName: filename,
        Path: filePath,
    }, nil  
}
