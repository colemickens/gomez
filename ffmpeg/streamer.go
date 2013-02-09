package ffmpeg

import (
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

type Transcoder struct {
	ffmpegExeName string
	prober        *Prober
}

type OutputFormat int

const (
	FormatWebM OutputFormat = iota
	FormatFLV  OutputFormat = iota + 1
)

func NewTranscoder(ffprobeExeName, ffmpegExeName string) (*Transcoder, error) {
	if ffprobeExeName == "" {
		ffprobeExeName = "ffprobe"
	}

	prober, err := NewProber(ffprobeExeName)
	if err != nil {
		return nil, err
	}

	return NewTranscoderWithProber(ffmpegExeName, prober)
}

func NewTranscoderWithProber(ffmpegExeName string, prober *Prober) (*Transcoder, error) {
	if ffmpegExeName == "" {
		ffmpegExeName = "ffmpeg"
	}
	return &Transcoder{ffmpegExeName, prober}, nil
}

func (s *Transcoder) TranscodeFile(filename string, format OutputFormat, offset, width, height, bitrate int) (ffStdOut io.Reader, ffStdErr io.Reader, err error) {

	// ffmpeg
	// -ss %o
	// -i %s
	// -threads 0
	// -async 1
	// -b %bk
	// -s %wx%h
	// -ar 44100
	// -ac 2
	// -v 0
	// -f flv
	// -vcodec libx264
	// -preset superfast
	// -

	cmdStr := ""
	if offset != 0 {
		cmdStr += "-ss {o} -async 1 "
	}
	cmdStr += "-i {s}"
	cmdStr += " -threads 4"
	if format == FormatFLV {
		cmdStr += " -b {b}k"
		cmdStr += " -s {w}x{h}"
		cmdStr += " -ar 44100"
		cmdStr += " -ac 2"
		cmdStr += " -v 0"
		cmdStr += " -f flv"
		cmdStr += " -vcodec libx264"
		cmdStr += " -preset superfast"
	} else if format == FormatWebM {
		cmdStr += " -codec:v libvpx"
		cmdStr += " -quality realtime -cpu-used 0 -b:v 10k"
		cmdStr += " -qmin 10 -qmax 42"
		cmdStr += " -maxrate 500k -bufsize 1000k"
		cmdStr += " -codec:a libvorbis -b:a 128k"
		cmdStr += " -f webm"
	}

	cmdStr += " -"

	ffp, err := s.prober.ProbeFile(filename)
	if err != nil {
		panic(err) // TODO: fix
	}

	realwidth := ffp.VideoStreams()[0].Width
	realheight := ffp.VideoStreams()[0].Height
	realratio := float64(realwidth) / float64(realheight)

	match := func(width, height int, ratio float64) error {
		if width%2 == 1 {
			log.Println("reject for width", width)
			return fmt.Errorf("Reject dimensions: width")
		}
		if height%2 == 1 {
			log.Println("reject for height", height)
			return fmt.Errorf("Reject dimensions: width")
		}
		if ratio-float64(width)/float64(height) > .0001 {
			log.Println("reject for ratio constraint")
			return fmt.Errorf("Reject dimensions: width")
		}
		return nil
	}

	/*
		if width == 0 || height == 0 {
			width = realwidth
			height = realheight
		}

		width -= width % 2
		for match(width, height, realratio) != nil{
			width -= 2
			height = int((float64(width) * float64(realheight)) / float64(realwidth))
		}
	*/

	if err = match(width, height, realratio); err != nil {
		// send an error back, we can't do this.
		return nil, nil, err
	}

	cmdStr = strings.Replace(cmdStr, "{o}", fmt.Sprintf("%d", offset), -1)
	cmdStr = strings.Replace(cmdStr, "{b}", fmt.Sprintf("%d", bitrate), -1)
	cmdStr = strings.Replace(cmdStr, "{w}", fmt.Sprintf("%d", width), -1)
	cmdStr = strings.Replace(cmdStr, "{h}", fmt.Sprintf("%d", height), -1)
	cmdParts := strings.Split(cmdStr, " ")

	// insert the filename into the cmdParts
	for i, v := range cmdParts {
		if v == "{s}" {
			cmdParts2 := append(cmdParts[:i], ffp.Format.Filename)
			cmdParts2 = append(cmdParts2, cmdParts[i+1:]...)
			cmdParts = cmdParts2
			break // NOTE: This assumes '%s' only occurs once in the transcode string!
		}
	}

	cmd := exec.Command(s.ffmpegExeName, cmdParts...)
	log.Println(cmd.Args)

	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	/*
		go io.Copy(w, stdoutPipe)
		go io.Copy(os.Stdout, stderrPipe)
	*/

	log.Println("starting", cmd.Args)
	cmd.Start()
	go cmd.Wait()

	return stdoutPipe, stderrPipe, nil // TODO: here, lolwut
}

// TODO : move this into the handler 
// w.Header().Set("Content-Type", contentType)
