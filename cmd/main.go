package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jawee/twitch-recorder/internal/configuration"
)

type Stream struct {
    Id string `json:"id"`
    UserId string `json:"user_id"`
    UserLogin string `json:"user_login"`
    UserName string `json:"user_name"`
    GameId string `json:"game_id"`
    GameName string `json:"game_name"`
    Type string `json:"type"`
    Title string `json:"title"`
    ViewerCount json.Number `json:"viewer_count"`
    StartedAt time.Time `json:"started_at"`
    Language string `json:"language"`
    ThumbnailUrl string `json:"thumbnail_url"`
    TagIds []string `json:"tag_ids"`
    IsMature bool `json:"is_mature"`
}

type SearchStream struct {
    Data []Stream `json:"data"`
}

type User struct {
	ID                 string    `json:"id"`
    Login              string    `json:"login"`
    DisplayName        string    `json:"display_name"`
    Type               string    `json:"type"`
    BroadcasterType    string    `json:"broadcaster_type"`
    Description        string    `json:"description"`
    ProfileImageUrl    string    `json:"profile_image_url"`
    OfflineImageUrl    string    `json:"offline_image_url"`
    ViewCount          string    `json:"view_count"`
    Email              string    `json:"email"`
    CreatedAt          time.Time `json:"created_at"`
}

type SearchUsers struct {
	Users []User `json:"data"`
}
type SearchChannel struct {
    Data []Channel `json:"data"`
}
type Channel struct {
	BroadcasterLanguage          string      `json:"broadcaster_language"`
	BroadcasterLogin             string      `json:"broadcaster_login"`
	DisplayName                  string      `json:"display_name"`
	GameId                       json.Number `json:"game_id"`
    GameName                     string      `json:"game_name"`
	ID                           json.Number `json:"id"`
    IsLive                       bool        `json:"is_live"`
    TagsIds                      []string    `json:"tags_ids"`
    ThumbnailUrl                 string      `json:"thumbnail_url"`
    Title                        string      `json:"title"`
    StartedAt                    time.Time    `json:"started_at"`
}

type TokenResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	Scope        []string `json:"scope"`
}

func main() {
    configuration := configuration.New()
    if configuration == nil {
        log.Println("Configuration is nil")
        os.Exit(1)
    }
    clientId := configuration.ClientId
    clientSecret := configuration.ClientSecret

    for {
        bearerToken, err := getTwitchBearerToken(clientId, clientSecret)
        if err != nil {
            log.Println(err)
            os.Exit(1)
        }

        for _, streamer := range configuration.Streamers {
            processStreamer(streamer, clientId, bearerToken)
        }
        time.Sleep(time.Minute)
    }
}


func processStreamer(username string, clientId string, bearerToken string) {
    log.Println("Processing streamer: " + username)
    users, err := getUserInformation(username, clientId, bearerToken)

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
        streams, err := getStreamInformation(user.ID, clientId, bearerToken)

        if err != nil {
            log.Println(err)
            os.Exit(1)
        }

        // str, err := json.Marshal(streams)
        // log.Printf("%s\n", str)
        if len(streams.Data) > 0 {
            log.Printf("%s is live\n", user.DisplayName)
            filename := fmt.Sprintf("%s_%s.mkv", streams.Data[0].StartedAt.Format("20060102_130405"), streams.Data[0].Title)
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
    filenamePath := baseDirectory + "/" + username + "/" + filename

    _, err := os.Stat(filenamePath)
    if err == nil {
        log.Println("File already exists")
        return
    }

    if _, err := os.Stat(baseDirectory + "/" + username); os.IsNotExist(err) {
        os.Mkdir(baseDirectory + "/" + username, 0777)
    }
    cmd := exec.Command("streamlink", "twitch.tv/" + username, "best", "-o", filenamePath)

    log.Println("Running cmd")
    cmd.Stdout = os.Stdout
    err = cmd.Run()

    if err != nil {
        log.Println(err)
    }

}

func getTwitchBearerToken(clientId string, clientSecret string) (string, error) {

    url := "https://id.twitch.tv/oauth2/token?client_id=" + clientId + "&client_secret=" + clientSecret + "&grant_type=client_credentials&scope=channel:read:subscriptions"

    req, err := http.NewRequest("POST", url, nil)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return "", err
    }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    var tokenResponse TokenResponse
    err = decoder.Decode(&tokenResponse)

    return tokenResponse.AccessToken, nil 
}

func getUserInformation(userName string, clientId string, bearerToken string) (*SearchUsers, error) {
    url := "https://api.twitch.tv/helix/users?login=" + userName

    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer " + bearerToken)
    req.Header.Set("Client-Id", clientId)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    var users SearchUsers
    err = decoder.Decode(&users)

    return &users, nil
}

func getChannelInformation(broadcasterId string, clientId string, bearerToken string) (*SearchChannel, error) {
    url := "https://api.twitch.tv/helix/search/channels?query=" + broadcasterId;
    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer " + bearerToken)
    req.Header.Set("Client-Id", clientId)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer resp.Body.Close()

    decoder := json.NewDecoder(resp.Body)
    var channels SearchChannel
    err = decoder.Decode(&channels)

    return &channels, nil
}

func getStreamInformation(userId string, clientId string, bearerToken string) (*SearchStream, error) {
    url := "https://api.twitch.tv/helix/streams?user_id=" + userId
    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("Authorization", "Bearer " + bearerToken)
    req.Header.Set("Client-Id", clientId)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer resp.Body.Close()

    decoder := json.NewDecoder(resp.Body)
    var streams SearchStream
    err = decoder.Decode(&streams)

    return &streams, nil
}
