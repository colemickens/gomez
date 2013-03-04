//+build ignore

package config

/*
	"type":     "postgres",
	"connection_string": "user=gomez dbname=gomez password=gomez_+234hjk;h98jkjjnasdf,bn,c.bvhgipuhtprfbv sslmode=disable"
*/

const SampleConfigStr = `{
	"db": {
		"type":     "sqlite3",
		"connection_string": "./gomez.db",
		"drop_on_start": "true"
	},

	"tmdb": {
		"api_key": "00ce627bd2e3caf1991f1be7f02fe12c",
		"project": "gomediaserver"
	},

	"tvdb": {
		"api_key": "78DAA2D23BE41064"
	},

	"library": {
		"ignore_keywords": [
			".AppleDouble", "VOB", "VIDEO_TS"
		],
		"follow_symlinks": true,
		"movie_dirs": [
			"/media/data/media/movies"
		],
		"music_dirs": [
			"/media/data/media/music"
		],
		"tvshow_dirs": [
			"/media/data/media/tvshows"
		]
	},

	"ffmpeg": {
		"ffmpeg_binary":  "/home/cole/Code_ext/ffmpeg_workbench/ffmpeg/ffmpeg",
		"ffprobe_binary": "/home/cole/Code_ext/ffmpeg_workbench/ffmpeg/ffprobe"
	},

	"web": {
		"web_listener":         "0.0.0.0:80",
		"api_listener":         "0.0.0.0:9999",
		"hostname":             "server.mickens.us",
		"swagger_path":         "/apidocs",
		"swagger_api":          "/apidocs.json"
	}
}`
