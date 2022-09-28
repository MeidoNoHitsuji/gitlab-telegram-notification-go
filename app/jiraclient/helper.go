package jiraclient

import (
	"context"
	"fmt"
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

func DeleteWorklogRecord(c *jira.Client, issueID string, worklogID string) (*jira.Response, error) {
	apiEndpoint := fmt.Sprintf("rest/api/2/issue/%s/worklog/%s", issueID, worklogID)
	req, err := c.NewRequestWithContext(context.Background(), "DELETE", apiEndpoint, nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, nil)
	if err != nil {
		jerr := jira.NewJiraError(resp, err)
		return resp, jerr
	}

	return resp, nil
}
