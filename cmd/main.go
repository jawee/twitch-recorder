package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/jawee/twitch-recorder/internal/configuration"
	"github.com/jawee/twitch-recorder/internal/discordclient"
	"github.com/jawee/twitch-recorder/internal/postprocessor"
	"github.com/jawee/twitch-recorder/internal/processor"
	"github.com/jawee/twitch-recorder/internal/recorder"
	"github.com/jawee/twitch-recorder/internal/recordingtracker"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

func main() {
    log.SetOutput(&lumberjack.Logger{
        Filename:   "/logs/twitch-recorder.log",
        MaxSize:    1, // megabytes
        MaxBackups: 0,
        MaxAge:     0, //days
    })
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
    rt := recordingtracker.New()
    c := make(chan *recorder.RecordedFile)

    postproc := postprocessor.New(disc)

    go startPostProcessing(postproc, c)

    proc := processor.New(c, twitchClient, rec, rt)

    for {

        for _, streamer := range configuration.Streamers {
            proc.ProcessStreamer(streamer)
        }
        time.Sleep(time.Minute)
    }
}

func startPostProcessing(postproc postprocessor.PostProcessor, c chan *recorder.RecordedFile) {
    for v := range c {
        go postproc.Process(v)
    }
}
