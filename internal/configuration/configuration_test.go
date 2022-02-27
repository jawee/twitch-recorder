package configuration

import "testing"

type MockConfigurationProvider struct {
}

func (m *MockConfigurationProvider) GetConfigurationJson() ([]byte, error) {
    return []byte("{\"client-id\":\"asdfafijewora\", \"client-secret\":\"asdfafijewora\", \"streamers\":\"streamer1, streamer2\", \"webhook-id\":\"asdads\", \"webhook-token\": \"asdfasdf\" }"), nil
}

func TestNewConfiguration(t *testing.T) {

    provider := new(MockConfigurationProvider)

    configuration, err := New(provider)

    if err != nil {
        t.Errorf("Error creating configuration: %s", err)
    }
    if configuration == nil {
        t.Errorf("Configuration is nil")
    }

    if configuration.ClientId != "asdfafijewora" {
        t.Errorf("ClientId is not set correctly")
    }

    if configuration.ClientSecret != "asdfafijewora" {
        t.Errorf("ClientSecret is not set correctly")
    }

    if configuration.StreamersString != "streamer1, streamer2" {
        t.Errorf("StreamersString is not set correctly")
    }

    if configuration.WebhookId != "asdads" {
        t.Errorf("WebhookId is not set correctly")
    }

    if configuration.WebhookToken != "asdfasdf" {
        t.Errorf("WebhookToken is not set correctly")
    }

    for idx, streamer := range configuration.Streamers {
        if streamer != "streamer1" && idx == 0 {
            t.Errorf("Streamer is not set correctly")
        }

        if streamer != "streamer2" && idx == 1 {
            t.Errorf("Streamer is not set correctly")
        }
    }

}
