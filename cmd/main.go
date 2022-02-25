package main

import (
	"log"
	"os"
	"time"

	"github.com/jawee/twitch-recorder/internal/configuration"
	"github.com/jawee/twitch-recorder/internal/processor"
	"github.com/jawee/twitch-recorder/internal/recorder"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

func main() {
    configProvider := new(configuration.FileConfigurationProvider)
    configuration := configuration.New(configProvider)
    if configuration == nil {
        log.Println("Configuration is nil")
        os.Exit(1)
    }
    clientId := configuration.ClientId
    clientSecret := configuration.ClientSecret

    twitchClient := twitch_client.New(clientId, clientSecret)
    baseDirectory := "/inprogress"
    rec := recorder.New(baseDirectory)
    c := make(chan *recorder.RecordedFile)
    proc := processor.New(c, twitchClient, rec)

    for {

        for _, streamer := range configuration.Streamers {
            proc.ProcessStreamer(streamer)
        }
        time.Sleep(time.Minute)
    }
}
