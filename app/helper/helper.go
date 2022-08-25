package helper

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"strings"
)

func Unique(arr []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range arr {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

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

func Keys(arr map[string]string) []string {
	var res []string

	for k := range arr {
		res = append(res, k)
	}

	return res
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

func TitleFirst(s string) string {
	
	if len(s) > 0 {
		upperS := strings.ToUpper(s)
		lowerS := strings.ToLower(s)
		s = fmt.Sprintf("%s%s", string([]rune(upperS)[0]), string([]rune(lowerS)[1:]))
	}

	return s
}
