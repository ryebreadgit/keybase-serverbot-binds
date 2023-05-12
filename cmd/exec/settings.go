package main

import (
	"encoding/json"
	"os"
)

type SettingsStruct struct {
	KeybaseUsername string `json:"keybase-bot-username"`
	Users           []struct {
		Username    string   `json:"username"`
		Permissions []string `json:"permissions"`
	} `json:"users"`
}

func GetSettings() *SettingsStruct {

	var settings SettingsStruct

	os.Mkdir("./settings/", os.ModePerm)

	file, err := os.ReadFile("./settings/settings.json")
	if err != nil {
		f, err := os.Create("./settings/settings.json")
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
