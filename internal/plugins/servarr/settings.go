package servarr

import (
	"encoding/json"
	"os"
)

type CredentialsStruct struct {
	OMDbAPI   string `json:"omdb-api-key"`
	RadarrAPI string `json:"radarr-api-key"`
	SonarrAPI string `json:"sonarr-api-key"`
}

type WebAddressStruct struct {
	RadarrLink string `json:"radarr-url"`
	SonarrLink string `json:"sondarr-url"`
}

type SettingsStruct struct {
	KeybaseUsername string            `json:"keybase-username"`
	WebAddresses    WebAddressStruct  `json:"web-addresses"`
	Credentials     CredentialsStruct `json:"credentials"`
}

func GetSettings() *SettingsStruct {

	var settings SettingsStruct

	os.Mkdir("./settings/", os.ModePerm)

	file, err := os.ReadFile("./settings/servarr.json")
	if err != nil {
		f, err := os.Create("./settings/servarr.json")
		if err != nil {
			panic(err)
		}
		data, err := json.MarshalIndent(settings, "", "  ")
		if err != nil {
			panic(err)
		}
		_, err2 := f.WriteString(string(data))
		if err2 != nil {
			panic(err)
		}
		f.Close()
	} else {

		json.Unmarshal(file, &settings)
	}
	return &settings
}

var settings *SettingsStruct = GetSettings()
