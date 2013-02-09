package tmdb

type Movie struct {
	Id           int    `json:"id"`
	Title        string `json:"title"`
	Backdrop_url string
	Poster_url   string
	Backdrops    []Image
	Posters      []Image
}

type Image struct {
	File_path    string
	Width        int
	Height       int
	Iso_639_1    interface{}
	Aspect_ratio float64
}

func convertImgPaths(_imgs []tmdbMovieImage) []Image {
	imgs := make([]Image, len(_imgs))
	for i, _img := range _imgs {
		imgs[i] = Image(_img)
		imgs[i].File_path = ImgMirror + imgs[i].File_path
	}
	return imgs
}
