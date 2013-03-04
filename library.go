package main

import (
	"github.com/eaigner/hood"
	"github.com/howeyc/fsnotify"

	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Library struct {
	Watcher *fsnotify.Watcher
}

func NewLibrary() (*Library, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("couldn't get watcher", err)
		return nil, err
	}

	return &Library{
		Watcher: watcher,
	}, nil
}

func (l *Library) Run() (err error) {
	for _, p := range cfg.Library.MovieDirs {
		log.Println(p)
		go l.scanFilesystem(p, movieFileHandler)
	}

	for _, p := range cfg.Library.TvshowDirs {
		log.Println(p)
		go l.scanFilesystem(p, tvshowFileHandler)
	}

	/*
		for _, p := range l.Cfg.Library.MusicDirs {
			_ = p
			go l.scanFilesystem(p, musicFileHandler, exitChan, true)
			exitCount++
		}
	*/

	for {
		l.findMetadata()
		time.Sleep(10 * time.Second)
	}

	return
}

func (l *Library) Stop() {
	//l.Watcher.Close()
	// shut down the findMetadata() loop
}

func movieFileHandler(f *FileRecord) {
	log.Println("movieFileHandler", f)
}

func tvshowFileHandler(f *FileRecord) {
	log.Println("tvshowFileHandler", f)
}

func musicFileHandler(f *FileRecord) {
	log.Println("musicFileHandler", f)
}

func (l *Library) findMetadata() {
	// do shit here, find metadata
	// concurrently fire out to
	// -> TMDB
	// -> TVDB
	// -> IMDB
	// -> CDDB
	// -> OpenSubtitles -> Hashing

	findFfprobeRecords := func() {
		type res struct {
			FileRecordId          int
			FfprobeFormatRecordId int
		}
		var results []res
		hd.Where("ffprobe_format_record.id", "=", "null").
			Join(hood.InnerJoin, &FfprobeFormatRecord{}, "file.id", "ffprobe_format_record.file_id").
			Find(results)
	}

	panic(results)

	findTmdbRecords := func() {
		type res struct {
			FileRecordId      int
			TmdbMovieRecordId int
		}
		results := make([]res)
		hd.Where("tmdb.id", "=", "null").
			Join(hood.InnerJoin, &MovieRecord{}, "file.id", "movie.file_id").
			Join(hood.InnerJoin, &TmdbRecord{}, "movie.id", "tmdb_movie_record.movie_id").
			Find(results)
	}

	findTvdbRecords := func() {
		hd.Whrere("tvdb.id", "=", "null").Join(hood.InnerJoin, &TvDb)
	}
}

/*
func (l *Library) mdQueue(mdChan chan *Entry) {
	for e := range mdChan {
		if e.Type == "movie" {
			log.Printf("- looking up '%s' with ss: '%s'\n", e.Name, e.Search_string)
			movie := tmdbApi.GetMovieViaSearch(e.Search_string)
			e.Tmdb = movie
			_, err := hd.Save(e)
			// can we just add to it?
			// can we do a pubsub directly into mongo?
			// I'd love a json pub-sub sync.
			// server pushes, client pushes, they meet in the middle...?
			//

			if err != nil {
				log.Println("db mdQueue update err", err)
				// use $addToSet ?
			}
		}
	}
}

func (l *Library) probeQueue(entryChan chan *Entry) {
	for e := range entryChan {
		// check if it exists in the library?

//			entry := &ffprobe.Output{}
//			l.Db.Find(map[string]string{"path:": e.Path}).One(&e.Ffprobe)
//			if entry != nil {
//			log.Println("possible already exists: ", e.Path, entry)
//				// probably exists
//				//continue
//			}


		probeinfo, err := prober.ProbeFile(e.Path)
		if err != nil {
			log.Println("------", e.Path, err)
			continue
		}
		if probeinfo == nil {
			log.Println("probeinfo == nil")
			continue
		}

		e.FfprobeOutput = probeinfo

		_, err = hd.Save(e)
		if err != nil {
			log.Println("db probeQueue update err", err)
		}
	}
}
*/

func (l *Library) scanFilesystem(walkPath string, cb func(*FileRecord)) {
	filepath.Walk(walkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Println("errrr", walkPath, path)
			return err
		}

		if func() bool {
			for _, ss := range cfg.Library.IgnoreKeywords {
				if strings.Contains(path, ss) {
					//log.Println(path, "contains", ss)
					return true
				}
			}
			return false
		}() {
			log.Println("skipping")
			return nil
		}

		linfo2, err := os.Lstat(path)

		if cfg.Library.FollowSymlinks && linfo2.Mode()&os.ModeSymlink == os.ModeSymlink {
			l.scanFsInternal(path, lt, entryChan, mdChan, c)
		} else if isMediaFile(path) {
			if !gd.PathExists(path) {
				f := &File{
					Path: path,
				}
				_, err = hd.Save(entry)
				if err != nil {
					panic(err)
				}
			} else {
				log.Println("skip existing file", path)
			}
		}
		return nil
	})
}

const videoextensions = ".m4v .3gp .nsv .ts .ty .strm .rm .rmvb .m3u .ifo .mov .qt .divx .xvid .bivx .vob .nrg .img .iso .pva .wmv .asf .asx .ogm .m2v .avi .bin .dat .dvr-ms .mpg .mpeg .mp4 .mkv .avc .vp3 .svq3 .nuv .viv .fli .flv .rar .001 .wpl .zip"

func isMediaFile(path string) bool {
	log.Println("checking file", path)
	for _, v := range strings.Split(videoextensions, " ") {
		if filepath.Ext(path) == v {
			return true
		}
	}
	return false
}

// TODO: This, better. Probably with precompiled regex
// do that in init()
// TODO: Do this for... uh... tvshows??
func scrubVideoString(path string) string {
	// strip file endings (er, doesn't this strip path?)
	info, err := os.Stat(path)
	if err != nil {
		panic(err)
	}

	s := info.Name()
	for _, v := range strings.Split(videoextensions, " ") {
		if filepath.Ext(s) == v {
			s = s[:len(s)-(len(v))]
			log.Println(s)
		}
	}
	for _, v := range []string{
		"480p", "720p", "dvd", "1080p", "webdl", "rip",
		"brrip", "readnfo", "xvid", "BluRay", "nHD",
		"extended edition", "BRRip", "READNFO", "XViD-TDP", "x264-NhaNc3",
		"extended", "bluray", "x264-crossbow",
		"UK", "(ENG)",
		"DTS", "x264-ESiR",
		"x264-BLOW", "PublicHD",
		"x264", "DTS-HDChina",
		"unrated", "BR", "QMax",
		// "xvid-{%s}", "x264-{%s}",
	} {
		// case (in)sensitive?
		s = strings.Replace(s, v, "", -1)
	}
	s = strings.Replace(s, ".", " ", -1)
	s = strings.Replace(s, "_", " ", -1)

	s = strings.Replace(s, "lotr", "lord of the rings", -1)
	s = strings.Replace(s, "SW", "star wars", -1)
	s = strings.Replace(s, "dir cut", "", -1)

	return s
}
