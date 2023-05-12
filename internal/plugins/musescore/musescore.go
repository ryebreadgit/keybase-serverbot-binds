package musescore

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func setCompArr(mstr MusescoreStruct, content string) MusescoreStruct {

	var txtmode string
	var arranger string
	var composer string

	mstr.Composer.Type = "Person"
	mstr.Arranger.Type = "Person"

	data := strings.Split(content, " ")
	for _, txt := range data {
		txt = strings.TrimSpace(txt)
		switch txt {
		case "Composer:":
			txtmode = "composer"
			continue
		case "Arranger:":
			txtmode = "arranger"
			continue
		}
		if txtmode == "composer" {
			composer += txt + " "
		}
		if txtmode == "arranger" {
			arranger += txt + " "
		}
	}

	if composer == "" && arranger == "" {
		composer = content
		arranger = content
	}

	if strings.Contains(composer, "\n") {
		composer = strings.Split(composer, "\n")[0]
	}
	if strings.Contains(arranger, "\n") {
		arranger = strings.Split(arranger, "\n")[0]
	}

	mstr.Composer.Name = strings.TrimSpace(composer)
	mstr.Arranger.Name = strings.TrimSpace(arranger)

	return mstr
}

func getDetails(user string, scoreid string) (MusescoreStruct, error) {
	ret := MusescoreStruct{}
	mslink := fmt.Sprintf("https://musescore.com/user/%v/scores/%v", user, scoreid)
	res, err := http.Get(mslink)
	if err != nil {
		return ret, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return ret, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return ret, err
	}

	fmt.Println(doc.Find("body").Text())

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		content := s.Text()
		if strings.Contains(content, `"@type": "MusicComposition"`) {
			json.Unmarshal([]byte(content), &ret)
		}
	})

	if ret.Composer.Name == "" {
		return ret, fmt.Errorf("no composor found for %v", scoreid)
	}

	ret = setCompArr(ret, ret.Composer.Name)

	return ret, nil
}
