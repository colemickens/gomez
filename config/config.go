package config

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
)

func ParseFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return ParseReader(f)
}

func SampleConfig() (*Config, error) {
	return ParseReader(bytes.NewBufferString(SampleConfigStr))
}

func ParseReader(reader io.Reader) (*Config, error) {
	var c Config
	dec := json.NewDecoder(reader)
	err := dec.Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

type Config struct {
	Db      DbConfig      `json:"db"`
	Tmdb    TmdbConfig    `json:"tmdb"`
	Tvdb    TvdbConfig    `json:"tvdb"`
	Library LibraryConfig `json:"library"`
	Ffmpeg  FfmpegConfig  `json:"ffmpeg"`
	Web     WebConfig     `json:"webui"`
}
type DbConfig struct {
	Type             string `json:"type"`
	ConnectionString string `json:"connection_string"`
	DropOnStart      bool   `json:"drop_on_start,string"`
}
type TvdbConfig struct {
	ApiKey string `json:"api_key"`
}
type TmdbConfig struct {
	ApiKey  string `json:"api_key"`
	Project string `json:"project"`
}
type LibraryConfig struct {
	IgnoreKeywords []string `json:"ignore_keywords"`
	FollowSymlinks bool     `json:"follow_symlinks"`
	MovieDirs      []string `json:"movie_dirs"`
	MusicDirs      []string `json:"music_dirs"`
	TvshowDirs     []string `json:"tvshow_dirs"`
}
type FfmpegConfig struct {
	FfmpegBinary  string `json:"ffmpeg_binary"`
	FfprobeBinary string `json:"ffprobe_binary"`
}
type WebConfig struct {
	ApiListener string `json:"api_listener"`
	WebListener string `json:"web_listener"`
	Hostname    string `json:"hostname"`
	SwaggerHome string `json:"swagger_home"`
	SwaggerPath string `json:"swagger_path"`
	SwaggerApi  string `json:"swagger_api"`
}
