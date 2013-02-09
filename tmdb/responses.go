package tmdb

type tmdbConfig struct {
	Images struct {
		Backdrop_sizes []string `json:"backdrop_sizes"`
		BaseUrl        string   `json:"base_url"` // don't know if this works.
		Poster_sizes   []string `json:"poster_sizes"`
		Profile_sizes  []string `json:"profile_sizes"`
	} `json:"image"`
}

type tmdbMovieResultPage struct {
	Page    int `json:"page"`
	Results []struct {
		Adult          bool    `json:"adult"`
		Backdrop_path  string  `json:"backdrop_path"`
		Id             int     `json:"id"`
		Original_title string  `json:"original_title"`
		Release_date   string  `json:"release_date"`
		Poster_path    string  `json:"poster_path"`
		Popularity     float64 `json:"popularity"`
		Title          string  `json:"title"`
	} `json:"results"`
	Total_pages   int `json:"total_pages"`
	Total_results int `json:"total_results"`
}

type tmdbMovieImages struct {
	Id        int              `json:"id"`
	Backdrops []tmdbMovieImage `json:"backdrops"`
	Posters   []tmdbMovieImage `json:"posters"`
}

type tmdbMovieImage struct {
	File_path    string
	Width        int
	Height       int
	Iso_639_1    interface{}
	Aspect_ratio float64
}
