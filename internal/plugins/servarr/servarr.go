package servarr

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
	log "github.com/sirupsen/logrus"
)

func SortName(data string) string {
	if data[:3] == "The" {
		return data[3:] + ", The"
	}
	return data
}

func RadarrGetAll(imdbID string, provider string) ([]RadarrTitleStruct, error) {
	var lookupURL string
	client := req.C()

	lookupURL = fmt.Sprintf("%vapi/v3/movie?apiKey=%v", settings.WebAddresses.RadarrLink, settings.Credentials.RadarrAPI)

	resp, err := client.R().Get(lookupURL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	data := []RadarrTitleStruct{}
	json.Unmarshal(resp.Bytes(), &data)

	return data, nil

}

func SonarrGetAll(imdbID string, provider string) ([]SonarrTitleStruct, error) {
	var lookupURL string
	client := req.C()

	lookupURL = fmt.Sprintf("%vapi/series?apiKey=%v", settings.WebAddresses.SonarrLink, settings.Credentials.SonarrAPI)

	resp, err := client.R().Get(lookupURL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	data := []SonarrTitleStruct{}
	json.Unmarshal(resp.Bytes(), &data)

	return data, nil
}

func RadarrAdd(data RadarrTitleStruct, provider string) error {
	var lookupURL string
	client := req.C()

	t, err := time.Parse("2006-01-02T00:00:00Z", data.InCinemas)
	if err != nil {
		return err
	}

	folder := fmt.Sprintf("/movies/%v (%v)", SortName(data.Title), t.Year())

	data.Monitored = true
	data.QualityProfileID = 4
	data.FolderName = folder
	data.Path = folder
	data.ID = 0

	lookupURL = fmt.Sprintf("%vapi/v3/movie?apikey=%v", settings.WebAddresses.RadarrLink, settings.Credentials.RadarrAPI)

	resp, err := client.R().SetBodyJsonMarshal(data).Post(lookupURL)

	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		log.Infoln(resp.String())
		return fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	return nil

}

func SonarrAdd(data SonarrTitleStruct, provider string) error {
	var lookupURL string
	client := req.C()

	folder := fmt.Sprintf("/tv/%v (%v)", SortName(data.Title), data.Year)

	data.Monitored = true
	data.QualityProfileID = 4
	data.SeasonFolder = true
	data.Path = folder

	for i := range data.Seasons {
		if i != 0 { // Skip special season
			data.Seasons[i].Monitored = true
		}
	}

	lookupURL = fmt.Sprintf("%vapi/series?apikey=%v", settings.WebAddresses.SonarrLink, settings.Credentials.SonarrAPI)

	temp, _ := json.Marshal(&data)

	resp, err := client.R().SetBodyJsonBytes(temp).Post(lookupURL)

	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("series non-200 status code: %d", resp.StatusCode)
	}

	newData := SonarrTitleStruct{}

	err = json.Unmarshal([]byte(resp.String()), &newData)
	if err != nil {
		return err
	}

	lookupURL = fmt.Sprintf("%vapi/v3/seasonPass?apikey=%v", settings.WebAddresses.SonarrLink, settings.Credentials.SonarrAPI)
	newBody := fmt.Sprintf(`{"series":[{"id":%d,"monitored":true}],"monitoringOptions":{"monitor":"all"}}`, newData.ID)

	time.Sleep(5 * time.Second) // Sleep for 5 seconds before changing monitoring status.

	resp, err = client.R().SetBodyString(newBody).Post(lookupURL)

	if err != nil {
		return err
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("seasonPass non-200 status code: %d", resp.StatusCode)
	}

	return nil

}

func ServarrLookup(imdbID string, provider string) (string, error) {
	var lookupURL string
	client := req.C()

	if provider == "Radarr" {
		lookupURL = fmt.Sprintf("%vapi/v3/movie/lookup?apiKey=%v&term=imdb:%v", settings.WebAddresses.RadarrLink, settings.Credentials.RadarrAPI, imdbID)
	} else {
		lookupURL = fmt.Sprintf("%vapi/series/lookup?apikey=%v&term=imdb:%v", settings.WebAddresses.SonarrLink, settings.Credentials.SonarrAPI, imdbID)
	}
	resp, err := client.R().Get(lookupURL)

	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", nil
	}

	return resp.String(), nil

}

func ServarrStartDownloads(provider string) error {
	var lookupURL string
	var cmd string
	client := req.C()

	if provider == "Radarr" {
		lookupURL = fmt.Sprintf("%vapi/v3/command?apiKey=%v", settings.WebAddresses.RadarrLink, settings.Credentials.RadarrAPI)
		cmd = `{"name": "MissingMoviesSearch"}`
	} else {
		lookupURL = fmt.Sprintf("%vapi/command?apikey=%v", settings.WebAddresses.SonarrLink, settings.Credentials.SonarrAPI)
		cmd = `{"name": "missingEpisodeSearch"}`
	}
	resp, err := client.R().SetBodyJsonString(cmd).Post(lookupURL)

	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		fmt.Println(resp.String())
		return fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	return nil

}
