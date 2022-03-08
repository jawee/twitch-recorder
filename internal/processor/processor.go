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
        return err
    }

    if users == nil || len(users.Users) == 0 {
        return errors.New("No users found for " + username)
    }

    for _, user := range users.Users {
        streams, err := sp.client.GetStreamInformation(user.ID)

        if err != nil {
            return err
        }

        if len(streams.Data) > 0 {
            log.Printf("%s is live\n", username)
            filename := fmt.Sprintf("%s_%s.mp4", streams.Data[0].StartedAt.Format("20060102_130405"), streams.Data[0].Title)
            filename = sanitizeFilename(filename)

            isRecording := sp.rt.IsAlreadyRecording(username)
            if isRecording {
                return fmt.Errorf("%s is already recording", username)
            }
            log.Printf("%s: Recording %s to %s\n", username, streams.Data[0].Title, filename)
            go func() {
                sp.rt.AddRecording(username)
                res, err := sp.rec.Record(username, filename)

                if err != nil {
                    log.Println("Recording error ", err)
                }

                if res != nil {
                    sp.c <- res
                }

                sp.rt.RemoveRecording(username)
            }()
        } else {
            return errors.New(user.DisplayName + " is offline")
        }
    }
    return nil
}

func sanitizeFilename(filename string) string {
    unallowedCharacters := []string {" ", "/", ":", "?", "&", "=", ",", "\"", "'", "\\", "*", "?", "!", "|", "<", ">", "#"}
    for _, chars := range unallowedCharacters {
        filename = strings.Replace(filename, chars, "_", -1)
    }
    return filename
}
