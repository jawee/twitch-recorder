package postprocessor

import "testing"

func TestNewFilenameNoSlash(t *testing.T) {
    expected := "/foo/bar/20220718_12345671_SomeStuff.mp4"
    fileName := "20220718_1234567_SomeStuff"
    filePath := "/foo/bar"
    res := getNewFilePath(filePath, fileName)

    if(res != expected) {
        t.Fatalf("Expected=%s, Got=%s", expected, res);
    }
}
func TestNewFilename(t *testing.T) {
    expected := "/foo/bar/20220718_12345671_SomeStuff.mp4"
    fileName := "20220718_1234567_SomeStuff"
    filePath := "/foo/bar/"
    res := getNewFilePath(filePath, fileName)

    if(res != expected) {
        t.Fatalf("Expected=%s, Got=%s", expected, res);
    }
}
