package config

import (
	"testing"
)

func checkStr(t *testing.T, left, right, msg string) {
	if left != right {
		t.Fatal(msg, "is wrong")
	}
}

func TestParseConfig(t *testing.T) {
	c, err := SampleConfig()
	if err != nil {
		t.Fatal("Parsing Test Config Failed", err)
	}

	checkStr(t, c.DbConfig.Type, "mongo", "config/db/type is wrong")
	checkStr(t, c.DbConfig.Hostname, "localhost", "config/db/hostname is wrong")
	checkStr(t, c.DbConfig.Dbname, "gomez", "config/db/dbname is wrong")
	checkStr(t, c.DbConfig.Username, "gomez", "config/db/username is wrong")
	checkStr(t, c.DbConfig.Password, "gomez", "config/db/password is wrong")

	checkStr(t, c.TmdbConfig.Project, "gomediaserver", "config/library/tmdb_project is wrong")
	checkStr(t, c.TmdbConfig.ApiKey, "00ce627bd2e3caf1991f1be7f02fe12c", "config/library/tmdb_api_key is wrong")

	checkStr(t, c.TvdbConfig.ApiKey, "78DAA2D23BE41064", "config/library/tvdb_api_key is wrong")

	if c.LibraryConfig.FollowSymlinks != true {
		t.Fatal("config/library/follow_symlinks is wrong")
	}

	checkStr(t, c.LibraryConfig.IgnoreKeywords[0], ".AppleDouble", "config/library/ignore_keywords is wrong")

	if len(c.LibraryConfig.TvshowDirs)+len(c.LibraryConfig.MovieDirs)+len(c.LibraryConfig.MusicDirs) != 3 {
		t.Error("config/library/{*}_dirs wrong length")
	}
	if c.LibraryConfig.TvshowDirs[0] != "/media/Videos/TV Shows/" ||
		c.LibraryConfig.MovieDirs[0] != "/media/Videos/Feature Length Films/" ||
		c.LibraryConfig.MusicDirs[0] != "/media/Music/" {
		t.Error("config/library/{*}_dirs wrong contents")
	}

	checkStr(t, c.WebuiConfig.Listener, "0.0.0.0:9000", "config/webui/listen interface is wrong")

	// TODO finish upnp test, not sure what I'm going to leave
	// changeable
}
