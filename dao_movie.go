package main

import (
	"github.com/eaigner/hood"
)

type TmdbMovie struct {
	// something embedded that gives me .Save() for this DAO?
	Id          hood.Id `json:"id"`
	FileId      int32   `json:"file_id"`
	TmdbId      int32   `json:"tmdb_id"`
	Title       string  `json:"title"`
	BackdropUrl string  `json:"backdrop_url"`
	PosterUrl   string  `json:"poster_url"`
}

type TmdbBackdrop struct {
	Id          hood.Id `json:"id"`
	TmdbMovieId int32   `json:"tmdb_movie_id"`
	Url         string  `json:"url"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
}

type TmdbPoster struct {
	Id          hood.Id `json:"id"`
	TmdbMovieId int32   `json:"tmdb_movie_id"`
	Url         string  `json:"url"`
	Width       int32   `json:"width"`
	Height      int32   `json:"height"`
}

type GenericDao struct {
	// inject hood hd or just use the global one?
}

//func (gd GenericDao) GetAllMovies(filters ...interface{}) ([]File, error) {
func (gd GenericDao) GetAllMovies() ([]File, error) {
	var movies []File
	err := hd.Find(&movies)
	if err != nil {
		panic(err)
	}
	return movies, nil
}

func (gd GenericDao) GetMovie(id int) (File, error) {
	var movies []File
	err := hd.Where("id = ?", id).Limit(1).Find(&movies)
	if err != nil {
		panic(err)
	}
	if len(movies) != 1 {
		return File{}, err
	}
	return movies[0], nil
}

func (gd GenericDao) AddNewFile(f *File) error {
	_, err := hd.Save(f)
	if err != nil {
		panic(err)
	}
	return nil
}

func (gd GenericDao) PathExists(path string) bool {
	var files []File
	err := hd.Where("path = ?", path).Limit(1).Find(&files)
	if err != nil {
		panic(err)
	}
	return (len(files) > 0)
}
