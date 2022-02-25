package processor

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jawee/twitch-recorder/internal/recorder"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

type Processor interface {
    ProcessStreamer(username string) error
}


type StreamProcessor struct {
    c chan *recorder.RecordedFile
    client twitch_client.InformationClient
    rec recorder.Recorder
}

func New(c chan *recorder.RecordedFile, client twitch_client.InformationClient, rec recorder.Recorder) *StreamProcessor {
    return &StreamProcessor{
        c: c,
        client: client,
        rec: rec,
    }
}

func (sp *StreamProcessor) ProcessStreamer(username string) error {
    log.Println("Processing streamer: " + username)

    users, err := sp.client.GetUserInformation(username)
    if err != nil {
        log.Println(err)
        return err
    }

    if users == nil || len(users.Users) == 0 {
        log.Println("No users found for " + username)
        log.Printf("users response: %v\n", users)
        return errors.New("No users found for " + username)
    }

    log.Println("Got user information for " + username)

    for _, user := range users.Users {
        log.Println("Getting stream information for " + user.Login)
        streams, err := sp.client.GetStreamInformation(user.ID)

        if err != nil {
            log.Println(err)
            return err
        }

        // str, err := json.Marshal(streams)
        // log.Printf("%s\n", str)
        if len(streams.Data) > 0 {
            log.Printf("%s is live\n", user.DisplayName)
            filename := fmt.Sprintf("%s_%s.mp4", streams.Data[0].StartedAt.Format("20060102_130405"), streams.Data[0].Title)
            filename = strings.Replace(filename, " ", "_", -1)

            log.Printf("Recording %s to %s\n", streams.Data[0].Title, filename)

            // baseDirectory := "/inprogress"
            // TODO this needs to be sent to a thread/goroutine. How to handle callback when done? 
            // c := make(chan *recorder.RecordedFile)
            go func() {
                // res, err := rec.Record(user.DisplayName, filename)
                res, err := sp.rec.Record(user.DisplayName, filename)
                if err != nil {
                    log.Println(err)
                }
                sp.c <- res
                close(sp.c)
            }()
        } else {
            log.Printf("%s is offline\n", user.DisplayName)
        }
    }
    return nil

}


