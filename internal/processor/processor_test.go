package processor

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jawee/twitch-recorder/internal/recorder"
	"github.com/jawee/twitch-recorder/internal/twitchclient"
)

type MockTwitchClient struct {
}

func (mtc *MockTwitchClient) GetChannelInformation(broadcasterId string) (*twitch_client.SearchChannel, error) {
    return &twitch_client.SearchChannel{
        Data: []twitch_client.Channel{
            {
                BroadcasterLanguage: "en",
                BroadcasterLogin: "somename",
                DisplayName: "Some Name",
                GameId: json.Number("12345"),
                GameName: "Some Game",
                ID: json.Number("12345"),
                IsLive: true,
                TagsIds: []string{"12345"},
                ThumbnailUrl: "https://static-cdn.jtvnw.net/jtv_user_pictures/somename-profile_image-1a1af1c2f8e7f9d4-300x300.png",
                Title: "Some Title",
                StartedAt: time.Now(),
            },
        }, 
    }, nil
}

func (mtc *MockTwitchClient) GetStreamInformation(userId string) (*twitch_client.SearchStream, error) {
    str := "2017-01-01T00:00:00Z"
    layout := "2006-01-02T15:04:05.000Z"
    t, _ := time.Parse(layout, str)
    
    return &twitch_client.SearchStream{
        Data: []twitch_client.Stream{
            {
                Id: "12345",
                UserId: "12345",
                UserLogin: "somename",
                UserName: "Some Name",
                GameId: "12345",
                GameName: "Some Game",
                Type: "live",
                Title: "Some Title",
                ViewerCount: json.Number("10000"),
                StartedAt: t,
                Language: "en",
                ThumbnailUrl: "https://static-cdn.jtvnw.net/jtv_user_pictures/somename-profile_image-1a1af1c2f8e7f9d4-300x300.png",
                TagIds: []string{"12345"},
                IsMature: false,
            },
        },
    }, nil
}

func (mtc *MockTwitchClient) GetUserInformation(userName string) (*twitch_client.SearchUsers, error) {
    // User with dummy data
    str := "2017-01-01T00:00:00Z"
    layout := "2006-01-02T15:04:05.000Z"
    t, _ := time.Parse(layout, str)

    return &twitch_client.SearchUsers{
        Users: []twitch_client.User{
            {
                ID: "12345",
                Login: "somename",
                DisplayName: "Some Name",
                Type: "staff",
                BroadcasterType: "Some bio",
                Description: "2017-01-01T00:00:00Z",
                ProfileImageUrl: "https://static-cdn.jtvnw.net/jtv_user_pictures/somename-profile_image-1a1af1c2f8e7f9d4-300x300.png",
                OfflineImageUrl: "https://static-cdn.jtvnw.net/jtv_user_pictures/somename-profile_banner-1a1af1c2f8e7f9d4-480.png",
                ViewCount: "10000",
                Email: "something@something.something",
                CreatedAt: t,
            },
        },
    }, nil
}   

type MockRecorder struct {
}

func (mr *MockRecorder) Record(username string, filename string) (*recorder.RecordedFile, error) {
    time.Sleep(time.Second * 5)
    return &recorder.RecordedFile{
        Name: filename + ".mp4",
        Path: "/some/path",
    }, nil
}


func TestProcessStreamerOnline(t *testing.T) {

    c := make(chan *recorder.RecordedFile)
    mockTwitchClient := new(MockTwitchClient)
    mockRecorder := new(MockRecorder)
    processor := New(&c, mockTwitchClient, mockRecorder)
    err := processor.ProcessStreamer("somename")

    if err != nil {
        t.Errorf("ProcessStreamer returned an error: %s", err)
    }

    res := <-c
    if res == nil {
        t.Errorf("ProcessStreamer did not return a result")
    }
}

func TestProcessStreamerOffline(t *testing.T) {

    c := make(chan *recorder.RecordedFile)
    mockTwitchClient := new(MockTwitchClient)
    mockRecorder := new(MockRecorder)
    processor := New(&c, mockTwitchClient, mockRecorder)
    err := processor.ProcessStreamer("somename")

    if err != nil {
        t.Errorf("ProcessStreamer returned an error: %s", err)
    }

    for msg := range c {
        if msg != nil {
            t.Errorf("ProcessStreamer returned a result: %s", msg)
        }
    }
}
