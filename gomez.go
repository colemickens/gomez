package main

import (
	// _ "github.com/bmizerany/pq"
	_ "github.com/mattn/go-sqlite3"

	_ "github.com/davecgh/go-spew/spew"

	"github.com/colemickens/gomez/config"
	"github.com/colemickens/gomez/ffmpeg"
	"github.com/colemickens/gomez/tmdb"
	"github.com/eaigner/hood"

	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	cfg        *config.Config
	db         *sql.DB
	hd         *hood.Hood
	tmdbApi    *tmdb.TmdbApi
	prober     *ffmpeg.Prober
	transcoder *ffmpeg.Transcoder
	library    *Library

	gd GenericDao
)

func init() {
	var err error

	// Parse Flags
	configFlag := flag.String("config", "", "configuration file")
	cmdFlag := flag.String("cmd", "", "command (writeconfig)")
	flag.Parse()

	// CMD: writeconfig
	if *cmdFlag == "writeconfig" {
		fmt.Fprintf(os.Stdout, "%s", config.SampleConfigStr)
		os.Exit(0)
	} else if *cmdFlag == "droptables" {
		// do that shit
	}

	// Read Config File
	if *configFlag != "" {
		cfg, err = config.ParseFile(*configFlag)
	} else {
		cfg, err = config.ParseFile("./config.json")
		if err != nil {
			cfg, err = config.SampleConfig()
		}
	}
	if err != nil {
		panic(err)
	}

	// Establish prober
	prober, err = ffmpeg.NewProber(cfg.Ffmpeg.FfprobeBinary)
	if err != nil {
		panic("main: couldn't initialize ffmpeg prober")
	}

	// Establish transcoder
	transcoder, err = ffmpeg.NewTranscoderWithProber(cfg.Ffmpeg.FfmpegBinary, prober)
	if err != nil {
		panic("main: couldn't initialize ffmpeg transcoder")
	}

	// Establish tmdb
	tmdbApi, err = tmdb.NewTmdbApi(cfg.Tmdb.ApiKey, nil)
	if err != nil {
		panic(err)
	}

	// Connect to Database
	hd, err = hood.Open(cfg.Db.Type, cfg.Db.ConnectionString)
	if err != nil {
		panic(err)
	}

	// Create Database Tables
	tables := []interface{}{
		&FileRecord{}, // dao_file.go

		&FfprobeFormatRecord{}, // dao_ffprobe.go
		&FfprobeStreamRecord{},

		&TmdbMovieRecord{}, // dao_movie.go
		&TmdbBackdropRecord{},
		&TmdbPosterRecord{},
	}
	for _, t := range tables {
		if cfg.Db.DropOnStart {
			err = hd.DropTableIfExists(t)
			if err != nil {
				panic(err)
			}
		}
		err = hd.CreateTable(t)
		if err != nil {
			panic(err)
		}
	}

	// Copy to a log file (this is acccessible via browser)
	lfile, err := os.Create("gomez.log")
	if err != nil {
		panic(err)
	}
	w := io.MultiWriter(lfile, os.Stderr)
	log.SetOutput(w)
}

func main() {
	var err error
	library, err = NewLibrary()
	if err != nil {
		panic(err)
	}
	go library.Run()

	// library.AddChannel("tvshows", {
	// 	rootPath := "/path/to/tvshows/"
	// });

	// TvChannelConfig?A

	// library.AddTvDir()
	// library.AddMovieDir()
	// library.AddDir(channel *Channel, )

	// webui.AddChannel("tvshows", TvShowChannel)

	// webui.AddChannel(Channel interface{whatver} )
	// 	l, err := library.NewLibrary(*cfg, tmdbApi, prober, db)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if cfg.Library.Enabled {
	// 		go l.Run()
	// 	}

	//upnpServer := upnpav.NewServer("Name", "something_else", handler)
	//go upnpServer.Run().

	go apiServices()
	webInterface()
}
