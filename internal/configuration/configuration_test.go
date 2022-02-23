package configuration

import "testing"

type MockConfigurationProvider struct {
}

func (m *MockConfigurationProvider) GetConfigurationJson() ([]byte, error) {
    return []byte("{\"client-id\":\"asdfafijewora\", \"client-secret\":\"asdfafijewora\", \"streamers\":\"streamer1, streamer2\" }"), nil
}

func TestNewConfiguration(t *testing.T) {

    provider := new(MockConfigurationProvider)

    configuration := New(provider)
    if configuration == nil {
        t.Errorf("Error creating configuration")
    }

    if configuration.ClientId != "asdfafijewora" {
        t.Errorf("Error creating configuration")
    }

    if configuration.ClientSecret != "asdfafijewora" {
        t.Errorf("Error creating configuration")
    }

    if configuration.StreamersString != "streamer1, streamer2" {
        t.Errorf("Error creating configuration")
    }

    for idx, streamer := range configuration.Streamers {
        if streamer != "streamer1" && idx == 0 {
            t.Errorf("Error creating configuration")
        }

        if streamer != "streamer2" && idx == 1 {
            t.Errorf("Error creating configuration")
        }
    }

}
