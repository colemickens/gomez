package main

import (
	"code.google.com/p/gorilla/mux"
	"github.com/colemickens/gomez/ffmpeg"

	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func webInterface() {
	var err error

	router := mux.NewRouter()
	r := router.Host(cfg.Web.Hostname).Subrouter()

	r.HandleFunc("/api/stream", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		entry := library.EntryForId(id)

		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		format := r.URL.Query().Get("format")
		width, _ := strconv.Atoi(r.URL.Query().Get("width"))
		height, _ := strconv.Atoi(r.URL.Query().Get("height"))
		bitrate, _ := strconv.Atoi(r.URL.Query().Get("bitrate"))

		filename := entry.Path
		destFormat := ffmpeg.FormatWebM

		if format == "flv" {
			destFormat = ffmpeg.FormatFLV
		}

		stdout, stderr, err := transcoder.TranscodeFile(filename, destFormat, offset, width, height, bitrate)
		if err != nil {
			panic(err) // TODO: fix
		}
		go io.Copy(os.Stdout, stderr)
		io.Copy(w, stdout)
	})

	list := func(w http.ResponseWriter, r *http.Request, t string) {
		w.Header().Set("Content-Type", "application/json")

		jsonEnc := json.NewEncoder(w)

		all := library.ListType(t)
		_ = jsonEnc.Encode(all)
	}

	r.HandleFunc("/api/movies", func(w http.ResponseWriter, r *http.Request) {
		list(w, r, "movie")
	})
	r.HandleFunc("/api/tvshows", func(w http.ResponseWriter, r *http.Request) {
		list(w, r, "tvshow")
	})
	r.HandleFunc("/api/music", func(w http.ResponseWriter, r *http.Request) {
		list(w, r, "music")
	})

	r.HandleFunc("/api/download", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		info := library.EntryForId(id)
		// filepath.EvalSymLinks?
		if r.URL.Query().Get("force") == "true" {
			w.Header().Set("Content-Disposition", "attachment")
		}
		http.ServeFile(w, r, info.Path)
	})

	r.HandleFunc("/api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")

		info := library.EntryForId(id)
		jsonEnc := json.NewEncoder(w)
		_ = jsonEnc.Encode(info)
	})

	r.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "gomez.log")
	})

	// serve files for the media server
	r.Handle("/{id:.*}", http.FileServer(http.Dir("./_webui/")))

	http.Handle("/", r)

	log.Println("Accepting (web) from:", cfg.Web.Hostname, "on:", cfg.Web.WebListener)
	err = http.ListenAndServe(cfg.Web.WebListener, nil)
	if err != nil {
		panic(err)
	}
}
