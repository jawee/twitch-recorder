package recorder

import (
	"log"
	"os"
	"os/exec"
	"path"
)

type RecordedFile struct {
    Name string
    Path string
}

type Recorder interface {
    Record(username string, filename string) (*RecordedFile, error)
}

type StreamlinkRecorder struct {
    baseDirectory string
}

func New(baseDirectory string) *StreamlinkRecorder {
    return &StreamlinkRecorder{
        baseDirectory: baseDirectory,
    }
}

func (s* StreamlinkRecorder) Record(username string, filename string) (*RecordedFile, error) {
    log.Println("Starting recording")

    filePath := path.Join(s.baseDirectory, username, filename)

    _, err := os.Stat(filePath)
    if err == nil {
        log.Println("File already exists")
        return nil, err
    }

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
        Name: filename,
        Path: filePath,
    }, nil  
}
