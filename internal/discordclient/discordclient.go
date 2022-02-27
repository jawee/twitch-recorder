package discordclient

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type webhookRequest struct {
    Content string `json:"content"`
}

type DiscordClient struct {
    webhookId string
    webhookToken string
}

func New(webhookId string, webhookToken string) *DiscordClient {
    return &DiscordClient{webhookId, webhookToken}
}

func (d *DiscordClient) SendMessage(message string) {
    webhookRequest := webhookRequest{Content: message}
    json, err := json.Marshal(webhookRequest)
    if err != nil {
        panic(err)
    }
    resp, err := http.Post("https://discordapp.com/api/webhooks/" + d.webhookId + "/" + d.webhookToken, "application/json", bytes.NewBuffer(json))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close() 
}
