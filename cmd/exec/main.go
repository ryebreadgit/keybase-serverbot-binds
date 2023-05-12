package main

import (
	"flag"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
	sblog "github.com/ryebreadgit/keybase-serverbot-binds/internal/logging"
	log "github.com/sirupsen/logrus"
)

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func commandNotFound(kbc *kbchat.API, msg *kbchat.SubscriptionMessage) {
	if _, err := kbc.SendMessage(msg.Message.Channel, "Command not found."); err != nil {
		log.Fatalln()
	}
	log.Infof("The following command from \"%v\" was not found: %v\n", msg.Message.Sender.Username, msg.Message.Content.Text.Body)
}

func init() {
	logTag := "ServerBot"
	sblog.LoggingInit(&logTag)
}

func main() {
	var kbc *kbchat.API
	var kbLoc string
	var err error

	log.Infoln("Starting serverbot keybase binds!")

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbchat.RunOptions{KeybaseLocation: kbLoc, StartService: true}); err != nil {
		log.Fatalf("Error creating Keybase API: %s\n", err.Error())
	}

	sub, err := kbc.ListenForNewTextMessages()
	if err != nil {
		log.Fatalf("Error listening on Keybase: %s\n", err.Error())
	}

	for {
		msg, err := sub.Read()
		if err != nil {
			log.Fatalf("failed to read message on Keybase: %s\n", err.Error())
		}

		user := msg.Message.Sender.Username

		var settings *SettingsStruct = GetSettings()

		if msg.Message.Content.TypeName != "text" || user == settings.KeybaseUsername { // Will ignore any non-text messages, messages sent by our bot, and non-valid users.
			continue
		}

		log.Debugf("Message recieved from \"%s\": %s", msg.Message.Sender.Username, msg.Message.Content.Text.Body)

		if !IsValidChannel(msg.Conversation.Channel.Name) {
			log.Debugf("The channel \"%s\" is not valid, skipping.", msg.Conversation.Channel.Name)
			continue
		}

		if bind := CheckBinds(msg); bind != "" {
			if CheckPermission(user, bind) {
				go RunBind(msg, bind)
			} else {
				commandNotFound(kbc, &msg)
			}
		} else {
			commandNotFound(kbc, &msg)
		}

	}
}
