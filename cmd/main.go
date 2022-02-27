package main

import (
	"log"
	"os"
	"time"

	"github.com/jawee/twitch-recorder/internal/configuration"
	"github.com/jawee/twitch-recorder/internal/discordclient"
	"github.com/jawee/twitch-recorder/internal/processor"
	"github.com/jawee/twitch-recorder/internal/recorder"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

func main() {
    configProvider := new(configuration.FileConfigurationProvider)
    configuration, err := configuration.New(configProvider)

    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
    if configuration == nil {
        log.Println("Configuration is nil")
        os.Exit(1)
    }
    clientId := configuration.ClientId
    clientSecret := configuration.ClientSecret
    discordToken := configuration.WebhookToken
    discordId := configuration.WebhookId

    twitchClient := twitch_client.New(clientId, clientSecret)
    baseDirectory := "/inprogress"
    disc := discordclient.New(discordId, discordToken)
    rec := recorder.New(baseDirectory, disc)
    c := make(chan *recorder.RecordedFile)
    proc := processor.New(c, twitchClient, rec)

    for {

        for _, streamer := range configuration.Streamers {
            proc.ProcessStreamer(streamer)
        }
        time.Sleep(time.Minute)
    }
}
