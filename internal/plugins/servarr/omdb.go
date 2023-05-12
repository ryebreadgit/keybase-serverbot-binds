package servarr

import (
	"encoding/json"
	"fmt"

	"github.com/imroc/req/v3"
)

func CheckOmdb(id string) (OmdbStruct, error) {

	var lookupURL string
	client := req.C()

	lookupURL = fmt.Sprintf("https://www.omdbapi.com/?i=%v&apikey=%v", id, settings.Credentials.OMDbAPI)

	resp, err := client.R().Get(lookupURL)

	if err != nil {
		return OmdbStruct{}, err
	}

	if resp.StatusCode != 200 {
		return OmdbStruct{}, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	data := OmdbStruct{}
	json.Unmarshal(resp.Bytes(), &data)

	return data, nil
}
