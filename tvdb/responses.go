package tvdb

import (
	"encoding/xml"
	"time"
)

type SeriesListResponse struct {
	XMLName xml.Name `xml:"Data"`
	Series  []Series
}

// <?xml version="1.0" encoding="UTF-8" ?>
// <Data>
//  <Series>
//   <seriesid>78107</seriesid>
//   <language>en</language>
//   <SeriesName>The Office (UK)</SeriesName>
//   <banner>graphical/78107-g9.jpg</banner>
//   <Overview>A mockumentary about life in a mid-sized suboffice paper merchants in a bleak British industrial town, where manager David Brent thinks he's the coolest, funniest, and most popular boss ever. He isn't. That doesn't stop him from embarrassing himself in front of the cameras on a regular basis, whether from his political sermonizing, his stand-up 'comedy', or his incredibly unique dancing. Meanwhile, long-suffering Tim longs after Dawn the engaged receptionist and keeps himself sane by playing childish practical jokes on his insufferable, army-obsessed deskmate Gareth. Will the Slough office be closed? Will the BBC give David a game show? Will Tim and Dawn end up with each other? And more importantly, will Gareth realize what a hopeless prat he is?</Overview>
//   <FirstAired>2001-07-01</FirstAired>
//   <IMDB_ID>tt0290978</IMDB_ID>
//   <id>78107</id>
//  </Series>
// </Data>

type Series struct {
	XMLName    xml.Name `xml:"Series"`
	Language   string   `xml:"language"`
	SeriesName string   `xml:"SeriesName"`
	BannerUrl  string   `xml:"banner"`
	Overview   string   `xml:"Overview"`
	FirstAired string   `xml:"FirstAired"`
	ImdbId     string   `xml:"IMDB_ID"`
	Id         int      `xml:"id"`
}

// <?xml version="1.0" encoding="UTF-8" ?>
// <Data>
//   <Episode>
//     <id>332179</id> 
//     <Combined_episodenumber>1</Combined_episodenumber> 
//     <Combined_season>1</Combined_season> 
//     <DVD_chapter /> 
//     <DVD_discid /> 
//     <DVD_episodenumber /> 
//     <DVD_Season /> 
//     <Director>McG</Director> 
//     <EpImgFlag /> 
//     <EpisodeName>Pilot</EpisodeName> 
//     <EpisodeNumber>1</EpisodeNumber> 
//     <FirstAired>2007-09-24</FirstAired> 
//     <GuestStars>Mieko Hillman|Kristine Blackport|Jim Pirri|Diana Gitelman|Mel Fair|Lynn A. Henderson|Odessa Rae|Jordan Potter|Tasha Campbell|Dale Dye|Matthew Bomer|Bruno Amato|Nicolas Pajon|Wendy Makkena</GuestStars> 
//     <IMDB_ID /> 
//     <language>en</language> 
//     <Overview>Chuck Bartowski is an average computer geek until files upon files of government secrets are downloaded into his brain. He is soon scouted by the CIA and NSA to act in place of their computer.</Overview> 
//     <ProductionCode /> 
//     <Rating /> 
//     <SeasonNumber>1</SeasonNumber> 
//     <Writer>Josh Schwartz|Chris Fedak</Writer> 
//     <absolute_number /> 
//     <filename>episodes/80348-332179.jpg</filename> 
//     <lastupdated>1209586232</lastupdated> 
//     <seasonid>27985</seasonid> 
//     <seriesid>80348</seriesid> 
//   </Episode>
// </Data>

type Episode struct {
	XMLName               xml.Name `xml:"Episode"`
	Id                    int
	CombinedEpisodeNumber int
	CombinedSeason        int
	Director              string
	EpisodeName           string
	EpisodeNumber         int
	FirstAired            string
	GuestStars            string // why is this not proper xml? jesus
	Language              string
	Overview              string
	SeasonNumber          int
	Writer                string
	Filename              string // wtf is this?
	LastUpdated           int64
	SeasonId              int
	SeriesId              int
}

// <?xml version="1.0" encoding="UTF-8" ?>
// <Actors>
//   <Actor>
//     <id>27747</id>
//     <Image>actors/27747.jpg</Image>
//     <Name>Matthew Fox</Name>
//     <Role>Jack Shephard</Role>
//     <SortOrder>0</SortOrder>
//   </Actor>
//   <Actor>
//     <id>27745</id>
//     <Image>actors/27745.jpg</Image>
//     <Name>Terry O'Quinn</Name>
//     <Role>John Locke</Role>
//     <SortOrder>0</SortOrder>
//   </Actor>
// </Actors>

type ActorsResponse struct {
	XMLName xml.Name `xml:"Actors"`
	Actors  []Actor
}

type Actor struct {
	XMLName   xml.Name `xml:"Actor"`
	Id        int      `xml:"id"`
	ImageUrl  string   `xml:"Image"`
	Name      string   `xml:"Name"`
	Role      string   `xml:"Role"`
	SortOrder int      `xml:"SortOrder"`
}

// <?xml version="1.0" encoding="UTF-8" ?>
// <Banners>
//   <Banner>
//     <id>89141</id>
//     <BannerPath>fanart/original/73739-34.jpg</BannerPath>
//     <BannerType>fanart</BannerType>
//     <BannerType2>1920x1080</BannerType2>
//     <Colors>|148,149,153|13,23,22|165,159,137|</Colors>
//     <Language>en</Language>
//     <Rating>7.6563</Rating>
//     <RatingCount>32</RatingCount>
//     <SeriesName>false</SeriesName>
//     <ThumbnailPath>_cache/fanart/original/73739-34.jpg</ThumbnailPath>
//     <VignettePath>fanart/vignette/73739-34.jpg</VignettePath>
//   </Banner>
// </Banners>
type BannersResponse struct {
	XMLName xml.Name `xml:"Banners"`
	Banners []Banner
}

type Banner struct {
	XMLName       xml.Name `xml:"Banner"`
	Id            int      `xml:"id"`
	BannerPath    string   `xml:"BannerPath"`
	BannerType    string   `xml:"BannerType"`
	BannerType2   string   `xml:"BannerType2"`
	Colors        string   `xml:"Colors"`
	Language      string   `xml:"Language"`
	Rating        float64  `xml:"Rating"`
	RatingCount   int      `xml:"RatingCount"`
	SeriesName    string   `xml:"SeriesName"`
	ThumbnailPath string   `xml:"ThumbailPath"`
	VignettePath  string   `xml:"VignettePath"`
}

// <?xml version="1.0" encoding="UTF-8" ?>
// <Data>
//   <Series>
//     <id>73739</id>
//     <Actors>|Matthew Fox|Terry O'Quinn|Evangeline Lilly|Naveen Andrews|Daniel Dae Kim|Yunjin Kim|Josh Holloway|Jorge Garcia|Elizabeth Mitchell|Henry Ian Cusick|Michael Emerson|Dominic Monaghan|Emilie de Ravin|Harold Perrineau Jr.|Ian Somerhalder|Maggie Grace|Malcolm David Kelley|John Terry|Andrew Divoff|Sam Anderson|M.C. Gainey|Zuleikha Robinson|L. Scott Caldwell|Nestor Carbonell|Kevin Durand|Jeff Fahey|Tania Raymonde|Mira Furlan|Alan Dale|Sonya Walger|Rebecca Mader|Ken Leung|Jeremy Davies|Kiele Sanchez|Rodrigo Santoro|Cynthia Watros|Adewale Akinnuoye-Agbaje|Michelle Rodriguez|</Actors>
//     <Airs_DayOfWeek>Tuesday</Airs_DayOfWeek>
//     <Airs_Time>9:00 PM</Airs_Time>
//     <ContentRating>TV-14</ContentRating>
//     <FirstAired>2004-09-22</FirstAired>
//     <Genre>|Action and Adventure|Drama|Science-Fiction|</Genre>
//     <IMDB_ID>tt0411008</IMDB_ID>
//     <Language>en</Language>
//     <Network>ABC</Network>
//     <NetworkID></NetworkID>
//     <Overview>After their plane, Oceanic Air flight 815, tore apart whilst thousands of miles off course, the survivors find themselves on a mysterious deserted island where they soon find out they are not alone.</Overview>
//     <Rating>9.1</Rating>
//     <RatingCount>637</RatingCount>
//     <Runtime>60</Runtime>
//     <SeriesID>24313</SeriesID>
//     <SeriesName>Lost</SeriesName>
//     <Status>Ended</Status>
//     <added></added>
//     <addedBy></addedBy>
//     <banner>graphical/73739-g4.jpg</banner>
//     <fanart>fanart/original/73739-34.jpg</fanart>
//     <lastupdated>1352667671</lastupdated>
//     <poster>posters/73739-11.jpg</poster>
//     <zap2it_id>SH672362</zap2it_id>
//   </Series>
//   <Episode>
//     <id>127151</id>
//     <Combined_episodenumber>1</Combined_episodenumber>
//     <Combined_season>0</Combined_season>
//     <DVD_chapter></DVD_chapter>
//     <DVD_discid></DVD_discid>
//     <DVD_episodenumber></DVD_episodenumber>
//     <DVD_season></DVD_season>
//     <Director></Director>
//     <EpImgFlag>1</EpImgFlag>
//     <EpisodeName>The Journey</EpisodeName>
//     <EpisodeNumber>1</EpisodeNumber>
//     <FirstAired>2005-04-27</FirstAired>
//     <GuestStars>|Brian Cox|</GuestStars>
//     <IMDB_ID></IMDB_ID>
//     <Language>en</Language>
//     <Overview>Flashbacks of the core characters illustrating who they were and what they were doing before the crash, a look at the island itself, and a preview of the big season finale.</Overview>
//     <ProductionCode>120</ProductionCode>
//     <Rating>7.5</Rating>
//     <RatingCount>4</RatingCount>
//     <SeasonNumber>0</SeasonNumber>
//     <Writer></Writer>
//     <absolute_number></absolute_number>
//     <airsafter_season></airsafter_season>
//     <airsbefore_episode>21</airsbefore_episode>
//     <airsbefore_season>1</airsbefore_season>
//     <filename>episodes/73739/127151.jpg</filename>
//     <lastupdated>1323264341</lastupdated>
//     <seasonid>21201</seasonid>
//     <seriesid>73739</seriesid>
//   </Episode>
//   [...]
// </Data>
type LangResponse struct {
	XMLName  xml.Name      `xml:"Data"`
	Series   LangSeries    `xml:"Series"`
	Episodes []LangEpisode `xml:"Episode"`
}

type LangSeries struct {
	Id            int    `xml:"id"`
	Actors        string `xml:"Actors"`
	AirDayOfWeek  string `xml:"Airs_DayOfWeek`
	AirTime       string `xml:"Airs_Time`
	ContentRating string `xml:"ContentRating`
	FirstAired    string `xml:"FirstAired`
	Genre         string `xml:"Genre`
	ImdbId        string `xml:"IMDB_ID`
	Language      string `xml:"Language"`
	Network       string `xml:"Network`
	NetworkId     string `xml:"NetworkID`
	Overview      string `xml:"Overview`
	Rating        string `xml:"Rating`
	RatingCount   string `xml:"RatingCount`
	Runtime       string `xml:"Runtime`
	SeriesId      string `xml:"SeriesID`
	SeriesName    string `xml:"SeriesName`
	Status        string `xml:"Status`
	Added         string `xml:"added`
	AddedBy       string `xml:"addedBy`
	Banner        string `xml:"banner`
	FanArt        string `xml:"fanart`
	LastUpdated   string `xml:"lastupdated`
	Poster        string `xml:"poster"`
	Zap2ItId      string "zap2it_id"
}

type LangEpisode struct {
	id int `xml:"id"`

	CombinedEpisodeNumber int       `xml:"Combined_episodenumber"`
	CombinedSeason        int       `xml:"Combined_season"`
	DvdChapter            int       `xml:"DVD_chapter"`
	DvdDiscId             int       `xml:"DVD_discid"`
	DvdEpisodeNumber      int       `xml:"DVD_episodenumber"`
	DvdSeason             int       `xml:"DVD_season"`
	Director              int       `xml:"Director"`
	EpisodeImageFlag      int       `xml:"EpImgFlag"`
	EpisodeName           int       `xml:"EpisodeName"`
	EpisodeNumber         int       `xml:"EpisodeNumber"`
	FirstAired            int       `xml:"FirstAired"`
	GuestStars            int       `xml:"GuestStars"`
	ImdbId                int       `xml:"IMDB_ID"`
	Language              int       `xml:"Language"`
	Overview              int       `xml:"Overview"`
	ProductionCode        int       `xml:"ProductionCode"`
	Rating                int       `xml:"Rating"`
	RatingCount           int       `xml:"RatingCount"`
	SeasonNumber          int       `xml:"SeasonNumber"`
	Writer                int       `xml:"Writer"`
	Absolute_number       int       `xml:"absolute_number"`
	AirsAfterSeason       int       `xml:"airsafter_season"`
	AirsBeforeEpisode     int       `xml:"airsbefore_episode"`
	AirsBeforeSeason      int       `xml:"airsbefore_season"`
	Filename              int       `xml:"filename"`
	LastUpdated           time.Time `xml:"lastupdated"`
	SeasonId              int       `xml:"seasonid"`
	SeriesId              int       `xml:"seriesid"`
}
