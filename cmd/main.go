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

    bearerToken, err := getTwitchBearerToken(clientId, clientSecret)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }
    // log.Println("Bearer token: " + bearerToken) 


    for _, streamer := range configuration.Streamers {
        processStreamer(streamer, clientId, bearerToken)
    }
    // processStreamer("bashbunni", clientId, bearerToken)

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

            startRecording(user.DisplayName, filename, "/tempdir")
        } else {
            log.Printf("%s is offline\n", user.DisplayName)
        }
    }

}

func startRecording(username string, filename string, path string) {
    log.Println("Starting recording")

    if _, err := os.Stat(path + "/" + username); os.IsNotExist(err) {
        os.Mkdir(path + "/" + username, 0777)
    }
    cmd := exec.Command("streamlink", "twitch.tv/" + username, "best", "-o", path + "/"+ username + "/" + filename)

    log.Println("Running cmd")
    cmd.Stdout = os.Stdout
    err := cmd.Run()

    if err != nil {
        log.Println(err)
    }

}


// POST https://id.twitch.tv/oauth2/token
//     ?client_id=<your client ID>
//     &client_secret=<your client secret>
//     &grant_type=client_credentials
//     &scope=<space-separated list of scopes>
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

// curl -X GET 'https://api.twitch.tv/helix/users?id=141981764' \
// -H 'Authorization: Bearer cfabdegwdoklmawdzdo98xt2fo512y' \
// -H 'Client-Id: uo6dggojyb8d6soh92zknwmi5ej1q2'

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
// curl -X GET 'https://api.twitch.tv/helix/search/channels?query=loserfruit' \
// -H 'Authorization: Bearer 2gbdx6oar67tqtcmt49t3wpcgycthx' \
// -H 'Client-Id: wbmytr93xzw8zbg0p1izqyzzc5mbiz'
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

// curl -X GET
// 'https://api.twitch.tv/helix/streams?user_id=12313123 \
// -H 'Authorization: Bearer 2gbdx6oar67tqtcmt49t3wpcgycthx' \
// -H 'Client-Id: uo6dggojyb8d6soh92zknwmi5ej1q2'

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
