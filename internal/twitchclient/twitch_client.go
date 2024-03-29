package twitch_client

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type InformationClient interface {
    GetChannelInformation(broadcasterId string) (*SearchChannel, error) 
    GetStreamInformation(userId string) (*SearchStream, error)
    GetUserInformation(userName string) (*SearchUsers, error)
}

type TwitchClient struct {
    clientSecret string
    clientId string
    bearerToken string
}

func New(clientId string, clientSecret string) *TwitchClient {
    return &TwitchClient{
        clientSecret: clientSecret,
        clientId: clientId,
    }
}

func (c *TwitchClient) GetChannelInformation(broadcasterId string) (*SearchChannel, error) {
    url := "https://api.twitch.tv/helix/search/channels?query=" + broadcasterId;

    var searchChannel SearchChannel
    err := c.makeAuthorizedRequest(url, &searchChannel, false)

    if err != nil {
        return nil, err
    }

    return &searchChannel, nil
}

func (c *TwitchClient) GetStreamInformation(userId string) (*SearchStream, error) {
    url := "https://api.twitch.tv/helix/streams?user_id=" + userId

    var searchStream SearchStream
    err := c.makeAuthorizedRequest(url, &searchStream, false)
    if err != nil {
        return nil, err
    }
    return &searchStream, nil
}

func (c *TwitchClient) GetUserInformation(userName string) (*SearchUsers, error) {
    url := "https://api.twitch.tv/helix/users?login=" + userName

    var searchUsers SearchUsers

    log.Printf("Making authorized request\n")
    err := c.makeAuthorizedRequest(url, &searchUsers, false)
    if err != nil {
        return nil, err
    }

    return &searchUsers, nil
}

func (c *TwitchClient) getTwitchBearerToken() error {
    url := "https://id.twitch.tv/oauth2/token?client_id=" + c.clientId + "&client_secret=" + c.clientSecret + "&grant_type=client_credentials&scope=channel:read:subscriptions"

    req, err := http.NewRequest("POST", url, nil)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return err
    }
    defer resp.Body.Close()
    decoder := json.NewDecoder(resp.Body)
    var tokenResponse TokenResponse
    err = decoder.Decode(&tokenResponse)

    if(err != nil) {
        log.Println(err)
        return err
    }

    c.bearerToken = tokenResponse.AccessToken

    return nil 
}

func (c *TwitchClient) makeAuthorizedRequest(url string, result interface{}, retry bool) error {
    log.Printf("makeAuthorizedRequest to: %s. Retry = %v\n", url, retry)
    if retry || c.bearerToken == "" {
        err := c.getTwitchBearerToken()
        if err != nil {
            log.Println(err)
            return err
        }
    }

    req, err := http.NewRequest("GET", url, nil)
    req.Header.Set("Client-Id", c.clientId)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.bearerToken)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println(err)
        return err
    }

    defer resp.Body.Close()

    if(resp.Status == "401 Unauthorized") {
        if !retry {
            c.makeAuthorizedRequest(url, result, true)
        } else {
            log.Println("401 Unauthorized")
            return errors.New(resp.Status)
        }
    }

    decoder := json.NewDecoder(resp.Body)
    err = decoder.Decode(result)
    if err != nil {
        log.Println(err)
        return err
    }

    return nil
}

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
    ViewCount          json.Number    `json:"view_count"`
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
