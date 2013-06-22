//+build ignore

package config

/*
	"type":     "postgres",
	"connection_string": "user=gomez dbname=gomez password=gomez_+234hjk;h98jkjjnasdf,bn,c.bvhgipuhtprfbv sslmode=disable"
*/

const SampleConfigStr = `{
	"db": {
		"type":              "sqlite3",
		"connection_string": "./gomez.db",
		"drop_on_start":     "true"
	},

	"scrapers": [
		"tmdb": {
			"api_key": "00ce627bd2e3caf1991f1be7f02fe12c",
			"project": "gomediaserver"
		},

		"tvdb": {
			"api_key": "78DAA2D23BE41064"
		}
	]

	"library": {
		"content": [
			{
				label: "movies",
				path: "/media/data/media/movies",
				scraper: "tmdb",
				type: "movies"
			},
			{
				label: "tvshows",
				path: "/media/data/media/tvshows",
				scraper: "tvdb",
				type: "tvshows"
			},
			{
				label: "documentaries",
				path: "/media/data/media/documentaries",
				scraper: "tvdb",
				type: "tvshows"
			},
			{
				label: "music",
				path: "/media/data/media/music",
				scraper: "cddb",
				type: "music"
			}
		],
		"settings": [
			"ignore_keywords": [
				".AppleDouble", "VOB", "VIDEO_TS"
			],
			"follow_symlinks": true,
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
