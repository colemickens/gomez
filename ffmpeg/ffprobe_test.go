package ffmpeg

import (
	"fmt"
	"testing"
)

func TestFfprobe(t *testing.T) {
	file := "/media/ExtendedVideos/Feature Length Films/Summer Storm/Cover/Summer_Storm_German-front.jpg"

	data, err := ProbeFile(file)
	if err != nil {
		t.Fatal("failed to probe", err)
	}

	fmt.Println(*data)
}
