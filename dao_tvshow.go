package main

import (
	"github.com/eaigner/hood"
)

type TvshowEpisodeFile struct {
	Id   hood.Id `json:"id"`
	Path string  `json:"path"`
}

type SearchStringToTvdbIdRecord struct {
	Id            hood.Id `json:"id"`
	SearchString  string  `json:"search_string"`   // search_string: "mike.and.molly"
	MatchedTvdbId int32   `json:"matched_tvdb_id"` // then go and check if we have SeriesRecord & EpisodeRecords
}

type TvdbSeries struct {
	Id     hood.Id
	TvdbId int32

	SeriesName string
}

type TvdbSeason struct {
	SeasonNumber
}

type TvdbEpisode struct {
	Id     hood.Id
	TvdbId int32

	SeasonNumber  int32
	EpisodeNumber int32
	EpisodeName   string

	TvdbSeasonId int32
	TvdbSeriesId int32
}
