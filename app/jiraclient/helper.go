package jiraclient

import "regexp"

func GetIssueKeyFromText(text string) string {
	re, _ := regexp.Compile(`(?m)^[A-Z-\d]+`)
	res := re.FindStringSubmatch(text)

	if len(res) == 0 {
		return ""
	}

	return res[0]
}
