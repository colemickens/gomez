package tvdb

// cache locally
// check for periodic updates, etc
//

import (
	"encoding/xml"
	//"log"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	mirror = "http://thetvdb.com" // their wiki says this can be hard-coded
)

type TvdbApi struct {
	ApiKey string
	Client *http.Client
}

func NewTvdbApi(apiKey string, client *http.Client) *TvdbApi {
	if client == nil {
		client = &http.Client{}
	}
	return &TvdbApi{apiKey, client}
}

func fetch(url string, obj interface{}) error {
	//log.Println("fetch:", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	//rdr := io.TeeReader(resp.Body, os.Stdout)
	//xmlDec := xml.NewDecoder(rdr)
	xmlDec := xml.NewDecoder(resp.Body)

	if err = xmlDec.Decode(obj); err != nil {
		return err
	}

	return nil
}

func (t *TvdbApi) GetExtendedInfo(seriesId int, language string) (br *BannersResponse, ar *ActorsResponse, ser *LangResponse, e error) {
	url := mirror + "/api/" + t.ApiKey + "/series/" + fmt.Sprintf("%d", seriesId) + "/all/" + language + ".zip"

	log.Println("GET " + url)

	buf := &bytes.Buffer{}
	resp, err := http.Get(url)
	io.Copy(buf, resp.Body)
	resp.Body.Close()

	reader := bytes.NewReader(buf.Bytes())
	r, err := zip.NewReader(reader, int64(reader.Len()))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to open zip reader for series id: " + fmt.Sprintf("%d", seriesId))
	}

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Println("error reading file in ExtendedInfo TVDB")
		}

		switch f.Name {
		case "banners.xml":
			//rdr := io.TeeReader(rc, os.Stdout)
			//xmlDec := xml.NewDecoder(rdr)
			xmlDec := xml.NewDecoder(rc)
			br = &BannersResponse{}
			err := xmlDec.Decode(br)
			if err != nil {
				log.Println("failed to decode banners.xml in ExtendedInfo TVDB")
			}
		case "actors.xml":
			xmlDec := xml.NewDecoder(rc)
			ar = &ActorsResponse{}
			err := xmlDec.Decode(ar)
			if err != nil {
				log.Println("failed to decode actors.xml in ExtendedInfo TVDB")
			}
		case language + ".xml":
			xmlDec := xml.NewDecoder(rc)
			ser = &LangResponse{}
			err := xmlDec.Decode(ser)
			if err != nil {
				log.Println("failed to decode lang.xml in ExtendedInfo TVDB")
			}
		default:
			log.Println("unknown file in ExtendedInfo TVDB")
		}
	}

	if br == nil {
		return nil, nil, nil, fmt.Errorf("banners.xml was missing for series id: %d", seriesId)
	}
	if ar == nil {
		return nil, nil, nil, fmt.Errorf("actors.xml was missing for series id: %d", seriesId)
	}
	if ser == nil {
		return nil, nil, nil, fmt.Errorf(language+".xml was missing for series id: %d", seriesId)
	}

	return
}

func (t *TvdbApi) GetSeries(seriesname, language string) (*SeriesListResponse, error) {
	url := mirror + "/api/GetSeries.php?seriesname=" + seriesname + "&language=" + language

	slr := &SeriesListResponse{}
	if err := fetch(url, slr); err != nil {
		return nil, err
	}

	return slr, nil
}

func (t *TvdbApi) GetSeriesByImdbId(imdbid string) (*SeriesListResponse, error) {
	url := mirror + "/api/GetSeriesByRemoteID.php?imdbid=" + imdbid // include language?

	sr := &SeriesListResponse{}
	if err := fetch(url, sr); err != nil {
		return nil, err
	}

	return sr, nil
}

/*func (t *TvdbApi) GetEpisodeByAirDate(seriesid, airdate, language string) (*Episode, error) {
	url := mirror + "/api/GetEpisodeByAirDate.php?apikey=" + t.ApiKey + "&seriesid=" + seriesid + "&airdate=" + airdate + "&language=" + language

	er := &EpisodeResponse{}
	if err := fetch(url, er); err != nil {
		return nil, err
	}

	return er, nil
}
*/
