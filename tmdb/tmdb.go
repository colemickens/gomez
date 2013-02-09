package tmdb

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	Mirror = "http://api.themoviedb.org/3" // their wiki says this can be hard-coded
	//Mirror = "http://private-5f5e-themoviedb.apiary.io/3"

	ImgMirror = "http://cf2.imgobject.com/t/p/w342"
)

// create an interface and make a cached version ??
// or like befor where we do it completel diff?

type TmdbApi struct {
	ApiKey string
	Config tmdbConfig
	Client *http.Client
}

func NewTmdbApi(apiKey string, client *http.Client) (*TmdbApi, error) {
	if client == nil {
		client = &http.Client{}
	}
	tmdbApi := &TmdbApi{
		ApiKey: apiKey,
		Client: client,
	}

	_url := Mirror + "/configuration"
	err := tmdbApi.get(&tmdbApi.Config, _url, url.Values{})
	if err != nil {
		return nil, err
	}

	return tmdbApi, nil
}

func (tmdbApi *TmdbApi) get(response interface{}, __url string, params url.Values) error {
	// add api key to the params
	// decode into response
	// TODO: seems like there's a better way to construct this url
	params.Set("api_key", tmdbApi.ApiKey)
	_url, err := url.Parse(__url + "?" + params.Encode())
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", _url.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")

	resp, err := tmdbApi.Client.Do(req)
	if err != nil {
		return err
	}

	jsonDec := json.NewDecoder(resp.Body)
	err = jsonDec.Decode(response)
	if err != nil {
		return err
	}

	resp.Body.Close()
	return nil
}

//
//
// Movies
func (t *TmdbApi) GetMovieInfo(id int) {
	_url := fmt.Sprintf("%s/movie/%d", Mirror, id)
	_ = _url
}
func (t *TmdbApi) GetMovieAltTitles(id int) {
	_url := fmt.Sprintf("%s/movie/%d/alternative_titles", Mirror, id)
	_ = _url
}
func (t *TmdbApi) GetMovieCastInfo(id int) {
	_url := fmt.Sprintf("%s/movie/%d/casts", Mirror, id)
	_ = _url
}
func (t *TmdbApi) GetMovieImages(id int) tmdbMovieImages {
	_url := fmt.Sprintf("%s/movie/%d/images", Mirror, id)

	movieImgs := tmdbMovieImages{}
	t.get(&movieImgs, _url, url.Values{})

	return movieImgs
}
func (t *TmdbApi) GetMoviePlotKeywords(id int) {
	_url := fmt.Sprintf("%s/movie/%d/keywords", Mirror, id)
	_ = _url
}
func (t *TmdbApi) GetMovieReleases(id int) {
	_url := fmt.Sprintf("%s/movie/%d/releases", Mirror, id)
	_ = _url
}
func (t *TmdbApi) GetMoveTrailers()        {}
func (t *TmdbApi) GetTranslations()        {}
func (t *TmdbApi) GetSimilarMovies()       {}
func (t *TmdbApi) GetBasicColectionInfo()  {}
func (t *TmdbApi) GetGetCollectionImages() {}

//
//
// People
func (t *TmdbApi) GetGetPersonInfo()    {}
func (t *TmdbApi) GetGetPersonCredits() {}
func (t *TmdbApi) GetGetPersonImages()  {}

//
//
// Company
func (t *TmdbApi) GetGetCompanyInfo()   {}
func (t *TmdbApi) GetGetCompanyMovies() {}

//
//
// Genres
func (t *TmdbApi) GetGetGenres()      {}
func (t *TmdbApi) GetGetGenreMovies() {}

//
//
// Search
func (t *TmdbApi) GetSearchMoviesByTitle(title string) tmdbMovieResultPage {
	_url := fmt.Sprintf("%s/search/movie", Mirror)
	values := url.Values{
		"query": []string{title},
	}

	movieResPage := tmdbMovieResultPage{}
	t.get(&movieResPage, _url, values)

	return movieResPage
}

func (t *TmdbApi) GetSearchMoviesByPerson(person string) {

}
func (t *TmdbApi) GetSearchCompaniesByName(name string) {

}

func (t *TmdbApi) GetMovieViaSearch(title string) *Movie {
	res := t.GetSearchMoviesByTitle(title)
	if res.Total_results == 0 {
		log.Printf("* failed to find with title: '%s'\n", title)
		/*
			return &Movie{
				Id:    -1,
				Title: "z",
			}
		*/
		return nil
	}

	res2 := t.GetMovieImages(res.Results[0].Id)

	// TODO: remove the hard coded image paths
	return &Movie{
		Id:           res.Results[0].Id,
		Title:        res.Results[0].Title,
		Backdrop_url: ImgMirror + res.Results[0].Backdrop_path,
		Poster_url:   ImgMirror + res.Results[0].Poster_path,
		Backdrops:    convertImgPaths(res2.Backdrops),
		Posters:      convertImgPaths(res2.Posters),
	}
}
