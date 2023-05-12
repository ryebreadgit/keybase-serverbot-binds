package main

import (
	"github.com/keybase/go-keybase-chat-bot/kbchat"
	"github.com/ryebreadgit/keybase-serverbot-binds/internal/plugins/musescore"
	"github.com/ryebreadgit/keybase-serverbot-binds/internal/plugins/servarr"
	log "github.com/sirupsen/logrus"
)

func CheckBinds(msg kbchat.SubscriptionMessage) string {
	if check, name := servarr.CheckMatch_(&msg); check {
		return name
	}
	if check, name := musescore.CheckMatch_(&msg); check {
		return name
	}
	return ""
}

func RunBind(msg kbchat.SubscriptionMessage, bind string) int {
	if bind == "" {
		log.Errorln("Error, attempted to run empty bind.")
		return -1
	}

	log.Infof("Running bind \"%v\"...", bind)

	switch bind {
	case "servarr":
		return servarr.Bind_(&msg)
	case "musescore":
		return musescore.Bind_(&msg)
	}
	log.Errorln("Error, bind \"%v\" not found.", bind)
	return -2
}
