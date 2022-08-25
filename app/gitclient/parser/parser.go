package parser

import (
	"gitlab-telegram-notification-go/helper"
	"regexp"
	"strings"
)

type ComType struct {
	Type        string
	Scope       string
	Exclamation bool
	Description string
	Body        string
	Footer      map[string][]string
}

func CompileCommit(s string) ComType {
	re, _ := regexp.Compile(`(?ms)(\w*)((?:\([^()\r\n]*\)|\()?(!)?)(:.*)`)
	res := re.FindAllStringSubmatch(s, -1)
	var c ComType
	for _, s := range res {
		for i, s2 := range s {
			switch i {
			case 1:
				c.Type = strings.ToLower(s2)
			case 2:
				c.Scope = strings.Trim(s2, "()")
			case 3:
				c.Exclamation = s2 != ""
			case 4:
				c.Description = strings.Trim(s2, " :")
			}
		}
	}

	if c.Type == "" {
		c.Type = "other"
	}

	if c.Scope == "" {
		c.Scope = "Другое"
	}

	if c.Description != "" {
		arr := strings.SplitN(c.Description, "\n\n", 2)
		if len(arr) > 1 {
			c.Description = strings.TrimSpace(arr[0])
			bodyFooter := arr[1]
			bodyFooter = strings.TrimSpace(strings.Trim(bodyFooter, "\n"))
			arr := strings.SplitN(bodyFooter, "\n\n", 2)
			if len(arr) == 2 {
				c.Body = strings.TrimSpace(arr[0])
				c.Footer = getFooter(arr[1])
			} else if len(arr) == 1 {
				matched, _ := regexp.MatchString(`(?m)(.*):(.*)`, arr[0])
				if matched {
					c.Footer = getFooter(arr[0])
					c.Body = ""
				} else {
					c.Body = arr[0]
					c.Footer = map[string][]string{}
				}
			} else {
				c.Body = ""
				c.Footer = map[string][]string{}
			}
		} else if len(arr) == 1 {
			c.Description = strings.TrimSpace(arr[0])
			c.Body = ""
			c.Footer = map[string][]string{}
		} else {
			c.Description = ""
			c.Body = ""
			c.Footer = map[string][]string{}
		}
	}

	c.Scope = helper.TitleFirst(c.Scope)
	c.Description = helper.TitleFirst(c.Description)
	c.Body = helper.TitleFirst(c.Body)

	return c
}

func getFooter(s string) map[string][]string {
	f := map[string][]string{}
	re, _ := regexp.Compile(`(?m)(.*):(.*)`)
	res := re.FindAllStringSubmatch(s, -1)

	for _, m := range res {
		if len(m) == 3 {
			k := strings.TrimSpace(m[1])
			v := strings.TrimSpace(m[2])
			q, ok := f[k]
			if ok {
				f[k] = append(q, strings.Split(v, ", ")...)
			} else {
				f[k] = strings.Split(v, ", ")
			}
		}
	}
	return f
}
