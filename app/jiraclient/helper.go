package jiraclient

import (
	"github.com/andygrunwald/go-jira"
	"regexp"
	"time"
)

func GetIssueKeyFromText(text string) string {
	re, _ := regexp.Compile(`(?m)^[A-Z-\d]+`)
	res := re.FindStringSubmatch(text)

	if len(res) == 0 {
		return ""
	}

	return res[0]
}

func getTime(original time.Time) *jira.Time {
	jiraTime := jira.Time(original)

	return &jiraTime
}
