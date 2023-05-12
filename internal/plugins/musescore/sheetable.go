package musescore

import (
	"fmt"
	"strings"

	"github.com/imroc/req/v3"
)

func getAuth(settings SettingsStruct) (string, error) {
	authUrl := fmt.Sprintf("%v/api/login", strings.TrimRight(settings.Sheetable.Url, "/"))
	client := req.C()
	resp, err := client.R().SetBodyJsonMarshal(settings.Sheetable.Credentials).Post(authUrl)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(resp.Status)
	}

	data := strings.Trim(resp.String(), "\"")

	return data, nil

}

func uploadSheet(settings SettingsStruct, auth string, sheetPath string, metadata MusescoreStruct) error {
	upUrl := fmt.Sprintf("%v/api/upload", strings.TrimRight(settings.Sheetable.Url, "/"))
	client := req.C()
	resp, err := client.R().SetBearerAuthToken(auth).SetFile("uploadFile", sheetPath).SetFormData(map[string]string{ // Set form data while uploading
		"sheetName":       metadata.Name,
		"composer":        metadata.Composer.Name,
		"releaseDate":     metadata.DatePublished[:strings.Index(metadata.DatePublished, "T")],
		"informationText": metadata.Text,
	}).Post(upUrl)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	if resp.StatusCode == 202 {
		return nil
	} else {
		return fmt.Errorf(resp.Status)
	}

}
