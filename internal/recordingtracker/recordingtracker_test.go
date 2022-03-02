package recordingtracker

import "testing"

func TestAdd(t *testing.T) {
    rt := New()
    rt.AddRecording("username")

    res := rt.IsAlreadyRecording("username")

    if res != true {
        t.Error("Expected true, got false")
    }

}

func TestNotRecording(t *testing.T) {
    rt := New()
    res := rt.IsAlreadyRecording("username")

    if res == true {
        t.Error("Expected false, got true")
    }
}

func TestAddTwo(t *testing.T) {
    rt := New()
    rt.AddRecording("username")
    rt.AddRecording("username2")

    res := rt.IsAlreadyRecording("username")

    if res != true {
        t.Error("Expected true, got false")
    }

    res = rt.IsAlreadyRecording("username2")

    if res != true {
        t.Error("Expected true, got false")
    }
}

func TestRemove(t *testing.T) {
    rt := New()
    rt.AddRecording("username")

    rt.RemoveRecording("username")
    res := rt.IsAlreadyRecording("username")

    if res == true {
        t.Error("Expected true, got false")
    }
}
