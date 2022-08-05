package helper

import (
	"github.com/xanzy/go-gitlab"
	"strings"
)

func Contains(arr []string, val string) bool {
	for _, a := range arr {
		if a == val {
			return true
		}
	}
	return false
}

func Slugify(val string) string {
	val = strings.Replace(strings.ToLower(val), "hook", "", 1)
	val = strings.TrimSpace(val)
	return strings.Replace(val, " ", "_", -1)
}

func AllowEvents() []string {
	events := []string{
		string(gitlab.EventTypePipeline),
		string(gitlab.EventTypeMergeRequest),
	}

	var newEvents []string

	for _, event := range events {
		e := Slugify(event)
		if len(e) != 0 {
			newEvents = append(newEvents, e)
		}

	}

	return newEvents
}
