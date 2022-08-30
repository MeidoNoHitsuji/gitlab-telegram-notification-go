package fm

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Link(text string, url string) string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", tgbotapi.EscapeText(tgbotapi.ModeHTML, url), tgbotapi.EscapeText(tgbotapi.ModeHTML, text))
}

func Underline(text string) string {
	return fmt.Sprintf("<u>%s</u>", text)
}

func Italic(text string) string {
	return fmt.Sprintf("<i>%s</i>", text)
}

func Bold(text string) string {
	return fmt.Sprintf("<b>%s</b>", text)
}
