package main

import (
	"github.com/eaigner/hood"
)

type TmdbMovieRecord struct {
	// something embedded that gives me .Save() for this DAO?
	Id          hood.Id `json:"id"`
	FileId      int32   `json:"file_id"`
	TmdbId      int32   `json:"tmdb_id"`
	Title       string  `json:"title"`
	BackdropUrl string  `json:"backdrop_url"`
	PosterUrl   string  `json:"poster_url"`
}

type TmdbBackdropRecord struct {
	Id          hood.Id `json:"id"`
	TmdbMovieId int32   `json:"tmdb_movie_id"`
	Url         string  `json:"url"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
}

type TmdbPosterRecord struct {
	Id          hood.Id `json:"id"`
	TmdbMovieId int32   `json:"tmdb_movie_id"`
	Url         string  `json:"url"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
}
