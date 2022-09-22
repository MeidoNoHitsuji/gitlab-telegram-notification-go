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

func Drop(arr []string, val string) []string {
	i := -1

	for j, s := range arr {
		if s == val {
			i = j
		}
	}

	if i != -1 {
		arr[i] = arr[len(arr)-1]
		arr[len(arr)-1] = ""
		arr = arr[:len(arr)-1]
	}

	return arr
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

func AllowEventsWithName() map[string]string {
	events := []string{
		string(gitlab.EventTypePipeline),
		string(gitlab.EventTypeMergeRequest),
	}

	newEvents := make(map[string]string)

	for _, event := range events {
		key := Slugify(event)
		if len(key) != 0 {
			newEvents[key] = event
		}
	}

	return newEvents
}

func AllowEvents() []string {
	return Keys(AllowEventsWithName())
}

func UpperFirst(s string) string {

	if len(s) > 0 {
		upperS := strings.ToUpper(s)
		s = fmt.Sprintf("%s%s", string([]rune(upperS)[0]), string([]rune(s)[1:]))
	}

	return s
}

func Grouping(parameters []string, maxElements int) [][]string {
	var result [][]string
	lines := len(parameters) / maxElements

	if len(parameters)%maxElements > 0 {
		lines++
	}

	for i := 0; i < lines; i++ {
		slice := (i + 1) * maxElements
		if slice > len(parameters) {
			slice = len(parameters)
		}
		result = append(result, parameters[i*maxElements:slice])
	}

	return result
}
