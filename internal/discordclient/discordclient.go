package discordclient

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type webhookRequest struct {
    Content string `json:"content"`
}

type NotificationClient interface {
    SendMessage(message string) error
}

type DiscordClient struct {
    webhookId string
    webhookToken string
}

func New(webhookId string, webhookToken string) *DiscordClient {
    log.Printf("Creating new discord client with webhook id %s and webhook token %s", webhookId, webhookToken)  
    return &DiscordClient{webhookId, webhookToken}
}

func (d *DiscordClient) SendMessage(message string) error {
    log.Println("Sending message to discord")
    webhookRequest := webhookRequest{Content: message}
    json, err := json.Marshal(webhookRequest)
    if err != nil {
        return err
    }
    uri := "https://discord.com/api/webhooks/" + d.webhookId + "/" + d.webhookToken
    resp, err := http.Post(uri, "application/json", bytes.NewBuffer(json))
    if err != nil {
        return err
    }
    defer resp.Body.Close() 

    log.Println(resp.Status)

    b, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    log.Println(string(b))

    return nil
}
