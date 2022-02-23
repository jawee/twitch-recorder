package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/jawee/twitch-recorder/internal/configuration"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

func main() {
    configuration := configuration.New()
    if configuration == nil {
        log.Println("Configuration is nil")
        os.Exit(1)
    }
    clientId := configuration.ClientId
    clientSecret := configuration.ClientSecret

    twitchClient := twitch_client.New(clientId, clientSecret)

    for {

        for _, streamer := range configuration.Streamers {
            processStreamer(streamer, twitchClient)
        }
        time.Sleep(time.Minute)
    }
}


func processStreamer(username string, twitchClient *twitch_client.TwitchClient) {
    log.Println("Processing streamer: " + username)

    users, err := twitchClient.GetUserInformation(username)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    if users == nil || len(users.Users) == 0 {
        log.Println("No users found for " + username)
        log.Printf("users response: %v\n", users)
        os.Exit(1)
    }

    log.Println("Got user information for " + username)

    for _, user := range users.Users {
        log.Println("Getting stream information for " + user.Login)
        streams, err := twitchClient.GetStreamInformation(user.ID)

        if err != nil {
            log.Println(err)
            os.Exit(1)
        }

        // str, err := json.Marshal(streams)
        // log.Printf("%s\n", str)
        if len(streams.Data) > 0 {
            log.Printf("%s is live\n", user.DisplayName)
            filename := fmt.Sprintf("%s_%s.mp4", streams.Data[0].StartedAt.Format("20060102_130405"), streams.Data[0].Title)
            filename = strings.Replace(filename, " ", "_", -1)

            log.Printf("Recording %s to %s\n", streams.Data[0].Title, filename)

            baseDirectory := "/inprogress"
            // TODO this needs to be sent to a thread/goroutine. How to handle callback when done? 
            go startRecording(user.DisplayName, filename, baseDirectory)
        } else {
            log.Printf("%s is offline\n", user.DisplayName)
        }
    }

}

func startRecording(username string, filename string, baseDirectory string) {
    log.Println("Starting recording")
    // filenamePath := baseDirectory + "/" + username + "/" + filename
    filePath := path.Join(baseDirectory, username, filename)

    _, err := os.Stat(filePath)
    if err == nil {
        log.Println("File already exists")
        return
    }

    userFolderPath := path.Join(baseDirectory, username)
    if _, err := os.Stat(userFolderPath); os.IsNotExist(err) {
        os.Mkdir(userFolderPath, 0777)
    }
    cmd := exec.Command("streamlink", "twitch.tv/" + username, "best", "-o", filePath)

    log.Println("Running cmd")
    cmd.Stdout = os.Stdout
    err = cmd.Run()

    if err != nil {
        log.Println(err)
    }

}

