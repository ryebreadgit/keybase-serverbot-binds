package main

import (
	"strings"

	log "github.com/sirupsen/logrus"
)

func getUserPermissions(username string) []string {
	for _, item := range settings.Users {
		if item.Username == username {
			return item.Permissions
		}
	}
	return nil
}

func CheckPermission(username string, permission string) bool {
	ret := false
	for _, item := range getUserPermissions(username) {
		switch item {
		case "all":
			ret = true
		case permission:
			ret = true
		}
	}
	return ret
}

func IsValidChannel(channel string) bool {
	s := strings.Split(channel, ",")
	if len(s) > 2 {
		log.Infof("The channel \"%v\" has more than 2 recipients, skipping as it not a direct message.", channel)
		return false
	}

	var validChannels []string

	for _, item := range settings.Users {
		validChannels = append(validChannels, item.Username)
	}

	for _, item := range s {
		if Contains(validChannels, item) {
			return true
		}
	}
	return false
}
