package processor

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jawee/twitch-recorder/internal/recorder"
	"github.com/jawee/twitch-recorder/internal/recordingtracker"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

type Processor interface {
    ProcessStreamer(username string) error
}


type StreamProcessor struct {
    c chan *recorder.RecordedFile
    client twitch_client.InformationClient
    rec recorder.Recorder
    rt *recordingtracker.RecordingTracker
}

func New(c chan *recorder.RecordedFile, client twitch_client.InformationClient, rec recorder.Recorder, rt *recordingtracker.RecordingTracker) *StreamProcessor {
    return &StreamProcessor{
        c: c,
        client: client,
        rec: rec,
        rt: rt,
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

        if len(streams.Data) > 0 {
            log.Printf("%s is live\n", user.DisplayName)
            filename := fmt.Sprintf("%s_%s.mp4", streams.Data[0].StartedAt.Format("20060102_130405"), streams.Data[0].Title)
            filename = strings.Replace(filename, " ", "_", -1)

            log.Printf("Recording %s to %s\n", streams.Data[0].Title, filename)
            go func() {

                isRecording := sp.rt.IsAlreadyRecording(username)
                if isRecording {
                    return
                }
                sp.rt.AddRecording(user.DisplayName)
                res, err := sp.rec.Record(user.DisplayName, filename)
                if err != nil {
                    log.Println(err)
                } else {
                    sp.c <- res
                }
                sp.rt.RemoveRecording(username)
                // Do I want a channel for each download, or just one?
                //close(sp.c)
            }()
        } else {
            log.Printf("%s is offline\n", user.DisplayName)
            return errors.New(user.DisplayName + " is offline")
        }
    }
    return nil

}


