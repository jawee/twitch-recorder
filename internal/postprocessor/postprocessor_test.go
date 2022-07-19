package postprocessor

import "testing"

func TestNewFilename(t *testing.T) {
    expected := "/foo/bar/20220718_12345671_SomeStuff.mp4"
    fileName := "20220718_1234567_SomeStuff.mp4"
    filePath := "/foo/bar/"
    res := getNewFilePath(filePath, fileName)

    if(res != expected) {
        t.Fatalf("Expected=%s, Got=%s", expected, res);
    }
}
