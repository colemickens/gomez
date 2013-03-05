package main

import (
	"github.com/eaigner/hood"
)

type MovieFile struct {
	Id   hood.Id `json:"id"`
	Path string  `json:"path"`
}

type MovieFileToTmdbMovie struct {
}

type TmdbMovie struct {
	Id     hood.Id `json:"id"`
	Title  string  `json:"title"`
	Year   int32   `json:"year"` //? 

	// do I need these vvvvvvvvvvvv
	BackdropUrl string `json:"backdrop_url"`
	PosterUrl   string `json:"poster_url"`
}

type TmdbBackdrop struct {
	Id          hood.Id `json:"id"`
	TmdbMovieId int32   `json:"tmdb_movie_record_id"`
	Url         string  `json:"url"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
}

type TmdbPoster struct {
	Id          hood.Id `json:"id"`
	TmdbMovieId int32   `json:"tmdb_movie_record_id"`
	Url         string  `json:"url"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
}

// path : provider
// /t/mm/mike.molly.s01e01 -> tvdb
// 
