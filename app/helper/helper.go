package helper

import (
	"github.com/xanzy/go-gitlab"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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
	var newS string

	for i := 0; i < len(s); i++ {
		if i == 0 {
			newS += cases.Upper(language.Und).String(string(s[i]))
		} else {
			newS += cases.Lower(language.Und).String(string(s[i]))
		}

	}

	return s
}
