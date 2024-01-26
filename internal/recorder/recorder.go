package recorder

import (
	"fmt"
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

    filePath := path.Join(s.baseDirectory, username, fmt.Sprintf("%s.mp4", filename))

    _, err := os.Stat(filePath)
    if err == nil {
        // call self with filename + "1"
        return s.Record(username, fmt.Sprintf("%s1", filename))
    }

    s.notificationClient.SendMessage("Starting recording for " + username + ". File " + fmt.Sprintf("%s.mp4", filename))

    userFolderPath := path.Join(s.baseDirectory, username)
    if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
        os.Mkdir(userFolderPath, 0777)
    }
    cmd := exec.Command("streamlink", "--twitch-disable-ads", "twitch.tv/" + username, "best", "-o", filePath)

    log.Println("Running cmd")
    cmd.Stdout = os.Stdout
    err = cmd.Run()

    return &RecordedFile{
        Username: username,
        FileName: filename,
        Path: filePath,
    }, err  
}
