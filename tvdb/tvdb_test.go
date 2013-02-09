package tvdb

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var tvdbapi = NewTvdbApi("78DAA2D23BE41064")

var (
	LostSeries = Series{
		Id:         73739,
		Language:   "en",
		SeriesName: "Lost",
		BannerUrl:  "graphical/73739-g4.jpg",
		Overview:   "After their plane, Oceanic Air flight 815, tore apart whilst thousands of miles off course, the survivors find themselves on a mysterious deserted island where they soon find out they are not alone.",
		FirstAired: "2004-09-22",
		ImdbId:     "tt0411008",
	}
	LostLangSeries = LangSeries{
		Id:            1,
		Actors:        "",
		AirDayOfWeek:  "",
		AirTime:       "",
		ContentRating: "",
		FirstAired:    "",
		Genre:         "",
		ImdbId:        "",
		Language:      "",
		Network:       "",
		NetworkId:     "",
		Overview:      "",
		Rating:        "",
		RatingCount:   "",
		Runtime:       "",
		SeriesId:      "",
		SeriesName:    "",
		Status:        "",
		Added:         "",
		AddedBy:       "",
		Banner:        "",
		FanArt:        "",
		LastUpdated:   "",
		Poster:        "",
		Zap2ItId:      "",
	}
	LostEpisode = Episode{
		Id: 127151,
		CombinedEpisodeNumber: 1,
		CombinedSeason:        0,
		Director:              "",
		EpisodeName:           "The Journey",
		EpisodeNumber:         1,
		FirstAired:            "2005-04-27",
		GuestStars:            "|Brian Cox|",
		Language:              "en",
		Overview:              "Flashbacks of the core characters illustrating who they were and what they were doing before the crash, a look at the island itself, and a preview of the big season finale.",
		SeasonNumber:          0,
		Writer:                "", // empty
		Filename:              "episodes/73739/127151.jpg",
		LastUpdated:           1323264341,
		SeasonId:              21201,
		SeriesId:              73739,
	}
	LostLangEpisode = LangEpisode{
		CombinedEpisodeNumber: 0,
		CombinedSeason:        0,
		DvdChapter:            0,
		DvdDiscId:             0,
		DvdEpisodeNumber:      0,
		DvdSeason:             0,
		Director:              0,
		EpisodeImageFlag:      0,
		EpisodeName:           0,
		EpisodeNumber:         0,
		FirstAired:            0,
		GuestStars:            0,
		ImdbId:                0,
		Language:              0,
		Overview:              0,
		ProductionCode:        0,
		Rating:                0,
		RatingCount:           0,
		SeasonNumber:          0,
		Writer:                0,
		Absolute_number:       0,
		AirsAfterSeason:       0,
		AirsBeforeEpisode:     0,
		AirsBeforeSeason:      0,
		Filename:              0,
		LastUpdated:           time.Now(),
		SeasonId:              0,
		SeriesId:              0,
	}
)

func TestGetSeries(t *testing.T) {
	slr, err := tvdbapi.GetSeries("LOST", "en")
	if err != nil {
		t.Fatal(err)
	}

	if len(slr.Series) != 87 {
		t.Fatal("\"LOST\" returned the wrong number of entries", 87, len(slr.Series))
	}

	if !reflect.DeepEqual(slr.Series[0], LostSeries) {
		t.Fatal("slr.Series[0] is incorrect")
	}
}

func TestGetSeriesByImdbId(t *testing.T) {
	slr, err := tvdbapi.GetSeriesByImdbId("tt0411008")
	if err != nil {
		t.Fatal(err)
	}

	if len(slr.Series) != 1 {
		t.Fatal("\"LOST\" returned the wrong number of entries")
	}

	if !reflect.DeepEqual(slr.Series[0], LostSeries) {
		t.Fatal("slr.Series[0] is incorrect")
	}
}

func TestGetExtendedInfo(t *testing.T) {
	seriesId := 73739
	language := "en"
	bannerResp, actorResp, langResp, err := tvdbapi.GetExtendedInfo(seriesId, language)
	if err != nil || bannerResp == nil || actorResp == nil || langResp == nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v\n", bannerResp)
	fmt.Printf("%#v\n", actorResp)
	fmt.Printf("%#v\n", langResp)

	// check banners
	if len(bannerResp.Banners) != 0 {
		t.Fatal("Wrong number of banners", 0, len(bannerResp.Banners))
	}

	// check actors
	if len(actorResp.Actors) != 0 {
		t.Fatal("Wrong number of actors", 0, len(actorResp.Actors))
	}

	if !reflect.DeepEqual(LostLangSeries, langResp.Series) {
		t.Fatal("Bad Lang Resp Series")
	}

	if len(langResp.Episodes) != 1 {
		t.Fatal("Wrong number of episodes", 1, len(langResp.Episodes))
	}

	if !reflect.DeepEqual(LostLangEpisode, langResp.Episodes[0]) {
		t.Fatal("Bad Lang Resp Series")
	}
	fmt.Println("__", langResp.Episodes[50], "__")
}
