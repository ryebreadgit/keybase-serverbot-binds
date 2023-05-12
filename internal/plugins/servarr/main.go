package servarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	sblog "github.com/ryebreadgit/keybase-serverbot-binds/internal/logging"
	log "github.com/sirupsen/logrus"
)

var kbc *kbchat.API

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func SendMsg(msg kbchat.SubscriptionMessage, body string) {
	if _, err := kbc.SendMessage(msg.Message.Channel, body); err != nil {
		log.Fatalln("Error sending the message \"%s\" due to the following error: %s", err.Error())
	}
	log.Infof("Message sent: %v\n", body)
}

func init() {
	logTag := "Servarr"
	sblog.LoggingInit(&logTag)
}

func main(msg kbchat.SubscriptionMessage) error {
	var err error
	if kbc, err = kbchat.Start(kbchat.RunOptions{}); err != nil {
		log.Fatalf("Error creating Servarr API: %s\n", err.Error())
	}

	SendMsg(msg, "IMDB Link found! Checking...")

	bod := msg.Message.Content.Text.Body

	imdbexp := regexp.MustCompile(`^(?:http:\/\/|https:\/\/)?(?:www\.)?(?:m\.)?(?:imdb.com\/title\/)?(tt[0-9]*)`)

	imdbgrp := imdbexp.FindStringSubmatch(bod)

	id := imdbgrp[1]

	omdbData, err := CheckOmdb(id)
	if err != nil {
		log.Errorf("Error, unable to complete OMDb lookup due to the following error: " + err.Error())
		SendMsg(msg, "Error, unable to complete. Please check logs.")
		return err
	}

	var provider string

	if omdbData.Type == "movie" {
		provider = "Radarr"
	} else if omdbData.Type == "series" {
		provider = "Sonarr"
	} else {
		log.Errorf("Error, unknown omdb type: " + omdbData.Type)
		SendMsg(msg, "Error, unable to complete. Please check logs.")
		return errors.New("unknown omdb type")
	}

	search, err := ServarrLookup(id, provider)
	if err != nil {
		log.Errorf("Error, unable to complete %v lookup due to the following error: %v", provider, err.Error())
		SendMsg(msg, "Error, unable to complete. Please check logs.")
		return err
	}

	if search == "" || search == "[]" { // If the response is blank that means we got a non-200 return code, let's check sonarr.
		retMsg := fmt.Sprintf("Error, no match found in %v for id \"%v\".", provider, id)
		log.Errorf(retMsg)
		SendMsg(msg, retMsg)
		return errors.New("no match found in " + provider)
	}

	var title string

	if provider == "Radarr" {

		data := []RadarrTitleStruct{}
		json.Unmarshal([]byte(search), &data)

		movies, _ := RadarrGetAll(id, provider)
		for _, item := range movies {
			if item.Title == data[0].Title {
				retMsg := fmt.Sprintf("The movie \"%v\" is already in Radarr, no changes made.", item.Title)
				log.Infof(retMsg)
				SendMsg(msg, retMsg)
				return nil
			}
		}

		title = data[0].Title
		err = RadarrAdd(data[0], provider)

	} else {
		data := []SonarrTitleStruct{}
		json.Unmarshal([]byte(search), &data)

		if data[0].Title == "IMDb at The Oscars" { // This is the default if our show is not found
			retMsg := fmt.Sprintf("The imdb id \"%v\" could not be found.", id)
			log.Errorln(retMsg)
			SendMsg(msg, retMsg)
			return errors.New("sonarr returned default data")
		}
		series, _ := SonarrGetAll(id, provider)
		for _, item := range series {
			if item.Title == data[0].Title {
				retMsg := fmt.Sprintf("The series \"%v\" is already in Sonarr, no changes made.", item.Title)
				log.Infof(retMsg)
				SendMsg(msg, retMsg)
				return nil
			}
		}

		title = data[0].Title
		err = SonarrAdd(data[0], provider)

	}

	if err != nil {
		log.Errorf("Error, unable to complete %v lookup due to the following error: %v", provider, err.Error())
		SendMsg(msg, "Error, unable to complete. Please check logs.")
		return err
	}

	err = ServarrStartDownloads(provider)

	if err != nil {
		log.Errorf("Error, unable to start downloads for %v due to the following error: %v", provider, err.Error())
		SendMsg(msg, "Error, unable to complete. Please check logs.")
		return err
	}

	SendMsg(msg, fmt.Sprintf("Successfully added \"%v\" to %v!", title, provider))

	return nil

}
