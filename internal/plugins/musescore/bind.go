package musescore

import (
	"regexp"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

var bindname string = "musescore"

// Bind_ is used for running the main function of the module returning the exit code.
func Bind_(msg *kbchat.SubscriptionMessage) int {
	err := main(*msg)
	if err != nil {
		return 1
	} else {
		return 0
	}
}

// Check match is used to see if the msg sent matches a specific criteria. If so, the Bind_ function will be ran using the msg data.
func CheckMatch_(msg *kbchat.SubscriptionMessage) (bool, string) {
	bod := msg.Message.Content.Text.Body
	imdbexp := regexp.MustCompile(`^(?:http:\/\/|https:\/\/)?(?:www\.)?(?:m\.)?(?:musescore.com\/user\/)?(?:[0-9]*)?(?:\/scores\/)?([0-9]*)`)

	imdbgrp := imdbexp.FindStringSubmatch(bod)

	if len(imdbgrp) == 0 { // If this is not an imdb link, ignore. Otherwise, let the user know we like what we see.
		return false, ""
	} else {
		return true, bindname
	}
}
