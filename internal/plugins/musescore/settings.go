package musescore

import (
	"encoding/json"
	"os"
)

type SettingsStruct struct {
	BaseDownloadPath string `json:"base_path"`
	Sheetable        struct {
		Url         string `json:"url"`
		Credentials struct {
			Username string `json:"email"`
			Password string `json:"password"`
		} `json:"credentials"`
	} `json:"sheetable"`
}

func GetSettings() *SettingsStruct {

	var settings SettingsStruct

	os.Mkdir("./settings/", os.ModePerm)

	file, err := os.ReadFile("./settings/musescore.json")
	if err != nil {
		f, err := os.Create("./settings/musescore.json")
		if err != nil {
			panic(err)
		}
		data, err := json.MarshalIndent(settings, "", "\t")
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
