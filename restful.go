package main

import (
	restful "github.com/emicklei/go-restful"

	"log"
	"strconv"
)

func getAllMovies(request *restful.Request, response *restful.Response) {
	var movies []TmdbMovieRecord
	hd.Find(movies)
	movies, err := gd.GetAllMovies()
	if err != nil {
		panic(err)
	}
	//entity := struct{ Movies []File }{movies}
	entity := movies
	log.Println(entity)
	response.WriteEntity(entity)
}
func getMovie(request *restful.Request, response *restful.Response) {
	id, err := strconv.Atoi(request.QueryParameter("id"))
	if err != nil {
		panic(err)
	}
	movie, err := gd.GetMovie(id)
	if err != nil {
		panic(err)
	}
	response.WriteEntity(movie)
}

func getAllAlbums(request *restful.Request, response *restful.Response)  {}
func getAllArtists(request *restful.Request, response *restful.Response) {}
func getAllSongs(request *restful.Request, response *restful.Response)   {}

func NewMovieService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/v2/movies").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("").To(getAllMovies).
		Doc("get all movies"))

	ws.Route(ws.GET("/{id}").To(getMovie).
		Doc("get a movie by its id").
		Param(ws.PathParameter("id", "movie id")))

	return ws
}

/*
func NewTvShowService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/v2/tvshows").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// / -> allEpisodes
	// /?show={showid}
	// /shows ?
	// /episodes ?
	ws.Route(ws.GET("/episodes").To(getAllEpisodes))

	ws.Route(ws.GET("/tvshows").To(getAllTvShows))

	ws.Route(ws.GET("/{id}").To(getMovie)).
		Doc("get a tvshow by its id").
		Param(ws.PathParameter("id", "movie id"))
}
*/

/*
func NewMusicService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/api/v2/music").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// ws.Route(ws.GET("").To(getAllMovies))

	// Albums
	ws.Route(ws.GET("/albums").To(getAllAlbums).
		Doc("get all albums"))

	// Artists
	ws.Route(ws.GET("/artists").To(getAllArtists).
		Doc("get all artists"))

	// Songs
	ws.Route(ws.GET("/songs").To(getAllSongs).
		Doc("get all songs"))

	ws.Route(ws.GET("/songs/{id}").
		Param(ws.PathParameter("id", "song id")))

	return ws
}
*/

func apiServices() {
	restful.Add(NewMovieService())
	//restful.Add(NewTvShowService())
	//restful.Add(NewMusicService())

	basePath := "http://" + cfg.Web.Hostname
	config := restful.SwaggerConfig{
		WebServicesUrl:  basePath,
		ApiPath:         cfg.Web.SwaggerApi,
		SwaggerPath:     cfg.Web.SwaggerPath,
		SwaggerFilePath: cfg.Web.SwaggerHome,
	}
	restful.InstallSwaggerService(config) // Add?

	/*
		log.Println("Accepting (api) from:", cfg.Web.Hostname, "on:", cfg.Web.ApiListener)
		err := http.ListenAndServe(cfg.Web.ApiListener, nil)
		if err != nil {
			panic(err)
		}
	*/
}
