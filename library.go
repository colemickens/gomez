package main

import (
	"github.com/colemickens/gomez/ffmpeg"
	"github.com/colemickens/gomez/tmdb"
	"github.com/howeyc/fsnotify"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type LibraryType int

const (
	MoviesLibType  LibraryType = iota
	TvshowsLibType LibraryType = iota
	MusicType      LibraryType = iota
)

type Library struct {
	Watcher *fsnotify.Watcher
}

type Entry struct {
	Id   int64  `json:"id" bson:"_id"`
	Name string `json:"name"`
	Path string `json:"path"`

	FfprobeOutput *ffmpeg.ProbeOutput `json:"ffprobe_output"`
	// ^^ Yuck, this is going to change
	// anyway as I move to RDBMS

	//TvdbEntry     *tvdb.Entry            `json:"tvdb"`
	//CddbEntry     *cddb.Entry            `json:"cddb"`

	Tmdb          *tmdb.Movie `json:"tmdb"`
	Type          string      `json:"type"`
	Search_string string      `json:"search_string"`
}

type ListEntry struct {
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Title         string  `json:"title"`
	Search_string string  `json:"search_string"`
	Duration      float64 `json:"duration"`
	Poster_url    string  `json:"poster_url"`
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

// should this be???
func (l *Library) Run() (err error) {
	go func() {
		for {
			select {
			case ev := <-l.Watcher.Event:
				log.Println(ev.String())
				if ev.IsDelete() {
					// TODO:
					// for shame, cole, write some tests
					// SHAME
				}
				if ev.IsRename() {
					// update the path in the db entry
					// delete the tmdb info and add it to the queue
					if err != nil {
						log.Println(err)
					}
				}
				if ev.IsCreate() {
					// ?
					// -- if file
					// -- -- ffprobe? (ffprobe queue?)
					// -- if dir
					// -- -- watch dir
					// -- -- ffprobe dir's contents (recursively)
					// -- -- should make scanFs add to watcher
					// -- -- -- then we just call scanFs(newDir)
					// -- -- -- it adds and probes. voila!
				}
			case err := <-l.Watcher.Error:
				log.Println("watcher err:", err)
			}
		}
	}()

	entryChan := make(chan *Entry)
	mdChan := make(chan *Entry)
	c := make(chan int)

	go l.probeQueue(entryChan)
	//go l.probeQueue(entryChan)
	//go l.probeQueue(entryChan)

	go l.mdQueue(mdChan)

	exitCount := 0
	for _, p := range cfg.Library.MovieDirs {
		log.Println(p)
		go l.scanFilesystem(p, "movie", entryChan, mdChan, c, true) // wait group these or what?
		exitCount++
	}
	/*
		for _, p := range l.Cfg.Library.MusicDirs {
			_ = p
			go l.scanFilesystem(p, "music", entryChan, mdChan, c, true)
			exitCount++
		}
		for _, p := range l.Cfg.Library.TvshowDirs {
			_ = p
			go l.scanFilesystem(p, "tvshow", entryChan, mdChan, c, true)
			exitCount++
		}
	*/

	for i := 0; i < exitCount; i++ {
		<-c
	}
	close(c)
	close(entryChan)

	return
}

func (l *Library) Stop() {
	l.Watcher.Close()
}

func (l *Library) PathForId(id string) string {
	var entry Entry
	hd.Where("id = ?", id).Find(&entry)
	log.Println(id)
	/*
		err := l.Db.FindId(bson.ObjectIdHex(id)).One(&entry)
		if err != nil {
			log.Println("FindId/One", err)
		}

		return entry.Path
	*/
	return ""
}

func (l *Library) ListType(t string) (listEntries []ListEntry) {
	entries := []Entry{}
	_ = hd.Find(&entries)

	for _, e := range entries {
		// TODO: Clean this up more
		/*
			if e.Tmdb == nil {
				log.Println("tmdb content is nil, can't procede")
				continue
			}
			if e.FfprobeOutput == nil {
				log.Println("ffprobe content is nil, can't procede")
				continue
			}
		*/
		var duration float64 = -1
		title := ""
		poster_url := ""

		if e.FfprobeOutput != nil {
			duration = e.FfprobeOutput.Format.Duration
		}
		if e.Tmdb != nil {
			title = e.Tmdb.Title
			poster_url = e.Tmdb.Poster_url
		}

		listEntries = append(listEntries, ListEntry{
			Id:            e.Id,
			Name:          e.Name,
			Title:         title,
			Search_string: e.Search_string,
			Duration:      duration,
			Poster_url:    poster_url,
		})
	}

	return
}

func (l *Library) EntryForId(id string) Entry {
	entry := Entry{}
	return entry
}

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
		/*
			entry := &ffprobe.Output{}
			l.Db.Find(map[string]string{"path:": e.Path}).One(&e.Ffprobe)
			if entry != nil {
			log.Println("possible already exists: ", e.Path, entry)
				// probably exists
				//continue
			}
		*/

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

// it comes in, we parse url
// put parsed url in
// queue loads from mongo db to link to synced db info
// also loads to attempt to link it to a tvdb/tmdb/cddb entry.

// TODO : man, how can I do this without toplevel?? :/
// or well, without using a separate function? I could extract the
// inner walk call, but I do alot around it and call this in
// several places :/
func (l *Library) scanFilesystem(walkPath string, lt string, entryChan chan *Entry, mdChan chan *Entry, c chan int, toplevel bool) {
	log.Println("before eval:", walkPath)
	var err error
	// TODO: Do I need to do this still?
	walkPath, err = filepath.EvalSymlinks(walkPath)
	if err != nil {
		log.Println("eval err", err)
	}
	log.Println("after eval:", walkPath)

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
			//log.Println("skipping")
			return nil
		}

		// TODO: add if l.Cfg.Follow_symlinks

		linfo2, err := os.Lstat(path)

		if linfo2.Mode()&os.ModeSymlink == os.ModeSymlink {
			l.scanFilesystem(path, lt, entryChan, mdChan, c, false)
		} else if info.IsDir() {
			if err := l.Watcher.Watch(path); err != nil {
				log.Println("err watching ", path, err)
			} else {
				log.Println("watching:", path)
			}
		} else if isMediaFile(path) {
			if !gd.PathExists(path) {
				f := File{
					Path:         path,
					SearchString: scrubVideoString(path), // do we need to store this even?
					// this is more of an accounting thing
				}
				gd.AddNewFile(&f)
				log.Println("added new file", f)
				// add it to a queue orrrrrr
				// let a worker pull from sql?
			} else {
				log.Println("skip existing file", path)
			}

			entry := &Entry{
				Name:          linfo2.Name(),
				Path:          path,
				FfprobeOutput: nil,
				Type:          lt,
				Search_string: scrubVideoString(path),
			}
			_, err = hd.Save(entry)
			if err != nil {
				log.Println("***", err)
			}
			entryChan <- entry
			mdChan <- entry
		}
		return nil
	})

	if toplevel {
		c <- 1
	}
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
