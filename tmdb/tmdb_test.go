package tmdb

import (
	"log"
	"testing"
)

func getTmdbApi(t *testing.T) *TmdbApi {
	tmdbApi, err := NewTmdbApi("00ce627bd2e3caf1991f1be7f02fe12c") // insert my test key, whatever or pipe it in?
	if err != nil {
		t.Fatal(err)
	}
	log.Println(tmdbApi.Config)
	return tmdbApi
}

func TestTmdbConfig(t *testing.T) {
	getTmdbApi(t)
}

func TestTmdbSearchMovies(t *testing.T) {
	tmdbApi := getTmdbApi(t)
	result := tmdbApi.GetSearchMoviesByTitle("prometheus") // ?
	log.Println(result.Results[0])
	// extract checkStr() checkInt() into a helper pkg, tired of redoing them repeatedly
}
