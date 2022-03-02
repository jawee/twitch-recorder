package recordingtracker

import (
	"log"
	"sync"
)

type RecordingTracker struct {
    mu sync.Mutex
    recordings []string
}

func New() *RecordingTracker {
    return &RecordingTracker{
        recordings: make([]string, 0),
    }
}

func (rt *RecordingTracker) IsAlreadyRecording(username string) bool {
    rt.mu.Lock()
    defer rt.mu.Unlock()

    for _, recording := range rt.recordings {
        if recording == username {
            log.Println(username + " is already recording")
            return true
        }
    }

    log.Println(username + " is not recording")
    return false
}


func (rt *RecordingTracker) AddRecording(username string) {
    log.Println("Adding recording for " + username)
    rt.mu.Lock()
    defer rt.mu.Unlock()

    rt.recordings = append(rt.recordings, username)
}

func (rt *RecordingTracker) RemoveRecording(username string) {
    log.Println("Removing recording for " + username)
    rt.mu.Lock()
    defer rt.mu.Unlock()

    for i, recording := range rt.recordings {
        if recording == username {
            rt.recordings = append(rt.recordings[:i], rt.recordings[i+1:]...)
            return
        }
    }
}
