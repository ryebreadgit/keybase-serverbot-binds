package servarr

import "time"

type RadarrTitleStruct struct {
	ID              int     `json:"id"`
	Title           string  `json:"title"`
	SortTitle       string  `json:"sortTitle"`
	SizeOnDisk      int     `json:"sizeOnDisk"`
	Overview        string  `json:"overview"`
	InCinemas       string  `json:"inCinemas"`
	PhysicalRelease *string `json:"physicalRelease"`
	Images          []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
		RemoteURL string `json:"remoteUrl"`
	} `json:"images"`
	Website             string   `json:"website"`
	Year                int      `json:"year"`
	HasFile             bool     `json:"hasFile"`
	YouTubeTrailerID    string   `json:"youTubeTrailerId"`
	Studio              string   `json:"studio"`
	Path                string   `json:"path"`
	RootFolderPath      string   `json:"rootFolderPath"`
	QualityProfileID    int      `json:"qualityProfileId"`
	Monitored           bool     `json:"monitored"`
	MinimumAvailability string   `json:"minimumAvailability"`
	IsAvailable         bool     `json:"isAvailable"`
	FolderName          string   `json:"folderName"`
	Runtime             int      `json:"runtime"`
	CleanTitle          string   `json:"cleanTitle"`
	ImdbID              string   `json:"imdbId"`
	TmdbID              int      `json:"tmdbId"`
	TitleSlug           string   `json:"titleSlug"`
	Certification       string   `json:"certification"`
	Genres              []string `json:"genres"`
	Tags                []int    `json:"tags"`
	Added               string   `json:"added"`
	Ratings             struct {
		Votes int `json:"votes"`
		Value int `json:"value"`
	} `json:"ratings"`
	Collection struct {
		Name   string `json:"name"`
		TmdbID int    `json:"tmdbId"`
		Images []struct {
			CoverType string `json:"coverType"`
			URL       string `json:"url"`
			RemoteURL string `json:"remoteUrl"`
		} `json:"images"`
	} `json:"collection"`
	Status string `json:"status"`
}

type SonarrTitleStruct struct {
	Title       string `json:"title"`
	Path        string `json:"path"`
	SortTitle   string `json:"sortTitle"`
	SeasonCount int    `json:"seasonCount"`
	Status      string `json:"status"`
	Overview    string `json:"overview"`
	Network     string `json:"network"`
	AirTime     string `json:"airTime"`
	Images      []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	RemotePoster string `json:"remotePoster"`
	Seasons      []struct {
		SeasonNumber int  `json:"seasonNumber"`
		Monitored    bool `json:"monitored"`
	} `json:"seasons"`
	Year              int       `json:"year"`
	ProfileID         int       `json:"profileId"`
	SeasonFolder      bool      `json:"seasonFolder"`
	Monitored         bool      `json:"monitored"`
	UseSceneNumbering bool      `json:"useSceneNumbering"`
	Runtime           int       `json:"runtime"`
	TvdbID            int       `json:"tvdbId"`
	TvRageID          int       `json:"tvRageId"`
	TvMazeID          int       `json:"tvMazeId"`
	FirstAired        time.Time `json:"firstAired"`
	SeriesType        string    `json:"seriesType"`
	CleanTitle        string    `json:"cleanTitle"`
	ImdbID            string    `json:"imdbId"`
	TitleSlug         string    `json:"titleSlug"`
	Certification     string    `json:"certification"`
	Genres            []string  `json:"genres"`
	Tags              []int     `json:"tags"`
	Added             time.Time `json:"added"`
	Ratings           struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	QualityProfileID int `json:"qualityProfileId"`
	ID               int `json:"id"`
}

type OmdbStruct struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Language string `json:"Language"`
	Country  string `json:"Country"`
	Awards   string `json:"Awards"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	ImdbRating string `json:"imdbRating"`
	ImdbVotes  string `json:"imdbVotes"`
	ImdbID     string `json:"imdbID"`
	Type       string `json:"Type"`
	Dvd        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
}
