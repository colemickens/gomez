package ffmpeg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
)

type Prober struct {
	ffprobeExeName string
}

type ProbeOutput struct {
	Format  Format   `json:"format"`
	Streams []Stream `json:"streams"`
}

type Format struct {
	Duration    float64           `json:"duration,string"`
	Filename    string            `json:"filename"`
	Format_name string            `json:"format_name"`
	Nb_sreams   int               `json:"nb_streams"`
	Size        int64             `json:"size,string"`
	Tags        map[string]string `json:"tags"`
}

type Stream struct {
	Codec_name string `json:"codec_name"`
	Codec_type string `json:"codec_type"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

func NewProber(ffprobeExeName string) (*Prober, error) {
	if ffprobeExeName == "" {
		ffprobeExeName = "ffmpeg"
	}

	// TODO: Check the version output!
	// ????? not even sure how to do this

	log.Println("CHECK PROBER VERSION", ffprobeExeName)

	return &Prober{ffprobeExeName}, nil
}

func (ffo *ProbeOutput) VideoStreams() (vstrms []Stream) {
	for _, s := range ffo.Streams {
		if s.Codec_type == "video" {
			vstrms = append(vstrms, s)
		}
	}
	return
}

func (ffo *ProbeOutput) AudioStreams() (astrms []Stream) {
	for _, s := range ffo.Streams {
		if s.Codec_type == "audio" {
			astrms = append(astrms, s)
		}
	}
	return
}

/*
{ "format" : { "bit_rate" : "2137889",
      "duration" : "87.336000",
      "filename" : "/home/cole/test1.mkv",
      "format_long_name" : "Matroska/WebM file format",
      "format_name" : "matroska,webm",
      "nb_streams" : 2,
      "size" : "23339337",
      "start_time" : "0.000000",
      "tags" : { "COMMENT" : "Matroska Validation File1, basic MPEG4.2 and MP3 with only SimpleBlock",
          "DATE_RELEASED" : "2010",
          "TITLE" : "Big Buck Bunny - test 1"
        }
    },
  "streams" : [ { "avg_frame_rate" : "1312499997/54687499",
        "codec_long_name" : "MPEG-4 part 2 Microsoft variant version 2",
        "codec_name" : "msmpeg4v2",
        "codec_tag" : "0x3234504d",
        "codec_tag_string" : "MP42",
        "codec_time_base" : "1/1000",
        "codec_type" : "video",
        "has_b_frames" : 0,
        "height" : 480,
        "index" : 0,
        "level" : -99,
        "pix_fmt" : "yuv420p",
        "r_frame_rate" : "24/1",
        "start_time" : "0.000000",
        "time_base" : "1/1000",
        "width" : 854
      },
      { "avg_frame_rate" : "125/3",
        "bits_per_sample" : 0,
        "channels" : 2,
        "codec_long_name" : "MP3 (MPEG audio layer 3)",
        "codec_name" : "mp3",
        "codec_tag" : "0x0000",
        "codec_tag_string" : "[0][0][0][0]",
        "codec_time_base" : "1/48000",
        "codec_type" : "audio",
        "index" : 1,
        "r_frame_rate" : "0/0",
        "sample_fmt" : "s16",
        "sample_rate" : "48000",
        "start_time" : "0.000000",
        "time_base" : "1/1000"
      }
    ]
}
*/

func (p *Prober) ProbeFile(path string) (*ProbeOutput, error) {
	// run command, check outputs
	bin := p.ffprobeExeName
	args := []string{
		"-v", "quiet",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
	}
	cmd := exec.Command(bin, append(args, path)...)
	fmt.Println(cmd.Path, cmd.Args)
	buf := &bytes.Buffer{}
	stdOut, _ := cmd.StdoutPipe()
	//stdErr, _ := cmd.StderrPipe()

	go io.Copy(buf, stdOut)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	cmd.Wait()

	data := &ProbeOutput{
	//Id: bson.NewObjectId(),
	}
	// TODO: ^ ?
	jsonDec := json.NewDecoder(buf)
	err = jsonDec.Decode(&data)
	if data.Format.Size == 0 {
		return nil, fmt.Errorf("data was empty")
	}
	if err != nil {
		return nil, err
	}

	_ = log.Println

	return data, nil
}

/*
{
    "streams": [
        {
            "index": 0,
            "codec_name": "flac",
            "codec_long_name": "FLAC (Free Lossless Audio Codec)",
            "codec_type": "audio",
            "codec_time_base": "1/44100",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "sample_fmt": "s16",
            "sample_rate": "44100",
            "channels": 2,
            "bits_per_sample": 0,
            "r_frame_rate": "0/0",
            "avg_frame_rate": "0/0",
            "time_base": "1/44100",
            "duration_ts": 7769244,
            "duration": "176.173333",
            "disposition": {
                "default": 0,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 0
            }
        },
        {
            "index": 1,
            "codec_name": "mjpeg",
            "codec_long_name": "MJPEG (Motion JPEG)",
            "codec_type": "video",
            "codec_time_base": "1/90000",
            "codec_tag_string": "[0][0][0][0]",
            "codec_tag": "0x0000",
            "width": 600,
            "height": 600,
            "has_b_frames": 0,
            "sample_aspect_ratio": "0:1",
            "display_aspect_ratio": "0:1",
            "pix_fmt": "yuvj422p",
            "level": -99,
            "r_frame_rate": "90000/1",
            "avg_frame_rate": "0/0",
            "time_base": "1/90000",
            "duration_ts": 15855600,
            "duration": "176.173333",
            "disposition": {
                "default": 0,
                "dub": 0,
                "original": 0,
                "comment": 0,
                "lyrics": 0,
                "karaoke": 0,
                "forced": 0,
                "hearing_impaired": 0,
                "visual_impaired": 0,
                "clean_effects": 0,
                "attached_pic": 1
            },
            "tags": {
                "comment": "Cover (front)"
            }
        }
    ],
    "format": {
        "filename": "04 - Don't Stop (Color On the Walls).flac",
        "nb_streams": 2,
        "format_name": "flac",
        "format_long_name": "raw FLAC",
        "duration": "176.173333",
        "size": "23412195",
        "bit_rate": "1063143",
        "tags": {
            "ARTIST": "Foster the People",
            "TITLE": "Don't Stop (Color On the Walls)",
            "ALBUM": "Torches [Best Buy Exclusive]",
            "DATE": "2011",
            "track": "04",
            "GENRE": "Indie"
        }
    }
}
*/
