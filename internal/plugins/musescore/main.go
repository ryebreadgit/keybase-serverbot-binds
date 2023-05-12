package musescore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/flytam/filenamify"
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
	logTag := "Musescore"
	sblog.LoggingInit(&logTag)
}

func writeMetadata(filepath string, msdata MusescoreStruct) error {

	b, err := json.MarshalIndent(msdata, "", "  ")

	if err != nil {
		return err
	}

	f, err := os.Create(filepath)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(string(b))

	if err != nil {
		return err
	}

	return nil

}

func dlMusescore(uri string, scorePath string, scoreTitle string) error {

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	// "/bin/bash", "-c", "npx",
	cmd := exec.Command("npx", "--yes", "dl-librescore@latest", "-i", uri, "-t", "pdf", "-o", scorePath)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	if err := cmd.Run(); err != nil {
		log.Error("Error, unable to startdl-librescore due to the following error: " + err.Error())
		return err
	}
	if stdout != nil && stdout.String() != "" {
		log.Debug(stdout)
	}
	if stderr != nil && stderr.String() != "" {
		if !strings.Contains(stderr.String(), "Done") {
			log.Error(stderr)
			return fmt.Errorf("stderr")
		}
		log.Debug(stderr)
	}
	return nil
}

func renamePdf(scorePath string, scoreTitle string) error {
	files, err := ioutil.ReadDir(scorePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".pdf" {
			filePath := fmt.Sprintf("%v/%v", scorePath, file.Name())
			newPath := fmt.Sprintf("%v/%v.pdf", scorePath, scoreTitle)
			os.Rename(filePath, newPath)
			return nil
		}
	}
	return fmt.Errorf("pdf not found")
}

func main(msg kbchat.SubscriptionMessage) error {
	var err error
	if kbc, err = kbchat.Start(kbchat.RunOptions{}); err != nil {
		log.Fatalf("Error creating Servarr API: %s\n", err.Error())
	}

	SendMsg(msg, "Musescore link found! Checking...")

	settings := *GetSettings()

	bod := msg.Message.Content.Text.Body

	msexp := regexp.MustCompile(`^(?:http:\/\/|https:\/\/)?(?:www\.)?(?:m\.)?(?:musescore.com\/user\/)?([0-9]*)?(?:\/scores\/)?([0-9]*)`)

	msgrp := msexp.FindStringSubmatch(bod)

	user := msgrp[1]
	scoreid := msgrp[2]

	det, err := getDetails(user, scoreid)
	if err != nil {
		SendMsg(msg, fmt.Sprintf("Error downloading score id \"%v\".", scoreid))
		return err
	}

	scoreTitle, _ := filenamify.Filenamify(det.Name, filenamify.Options{Replacement: "-"})
	scoreComposer, _ := filenamify.Filenamify(det.Composer.Name, filenamify.Options{Replacement: ""})
	scorePath := fmt.Sprintf("%v/%v/%v/", settings.BaseDownloadPath, scoreComposer, scoreTitle)
	scoreFile := fmt.Sprintf("%v/%v.pdf", scorePath, scoreTitle)

	if _, err := os.Stat(scoreFile); err == nil {
		SendMsg(msg, fmt.Sprintf("\"%v\" already downloaded, skipping.", scoreTitle))
		return nil
	}

	os.MkdirAll(scorePath, os.ModePerm)

	err = writeMetadata(fmt.Sprintf("%v/%v.json", scorePath, scoreTitle), det)
	if err != nil {
		SendMsg(msg, fmt.Sprintf("Error writing metadata \"%v\".", scoreid))
		log.Error(err)
		return err
	}

	err = dlMusescore(det.URL, scorePath, scoreTitle)
	if err != nil {
		SendMsg(msg, fmt.Sprintf("Error downloading \"%v\".", scoreid))
		log.Error(err)
		return err
	}

	err = renamePdf(scorePath, scoreTitle)
	if err != nil {
		SendMsg(msg, fmt.Sprintf("Error renaming \"%v\".", scoreid))
		log.Error(err)
		return err
	}

	sheetAuth, err := getAuth(settings)
	if err != nil {
		SendMsg(msg, fmt.Sprintf("Error getting sheetable auth for \"%v\".", scoreid))
		log.Error(err)
		return err
	}

	err = uploadSheet(settings, sheetAuth, scoreFile, det)
	if err != nil {
		SendMsg(msg, fmt.Sprintf("Error uploading \"%v\".", scoreid))
		log.Error(err)
		return err
	}

	SendMsg(msg, fmt.Sprintf("Successfully downloaded and added \"%v\"!", scoreTitle))

	return nil

}
