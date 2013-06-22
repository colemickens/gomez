package config

const SampleConfigStr = `{
	"db": {
		"type":              "postgres",
		"connection_string": "user=gomez dbname=gomez password=gomez234hjkh98jkjjnasdfb9088cbvhgipuhtprfbv sslmode=disable"
		"drop_on_start":     "true"
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
		"ffmpeg_binary":  "",
		"ffprobe_binary": ""
	},

	"web": {
		"web_listener":         "0.0.0.0:80",
		"api_listener":         "0.0.0.0:9999",
		"hostname":             "server.mickens.us",
		"swagger_path":         "/apidocs",
		"swagger_api":          "/apidocs.json"
	}
}`
