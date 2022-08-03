package command

import (
	"fmt"
	"github.com/xanzy/go-gitlab"
	"gitlab-telegram-notification-go/client"
	"log"
	"strings"
)

func Subscribe(Arguments string) (string, *gitlab.Project) {

	git := client.Gitlab()

	arguments := strings.Split(
		strings.TrimSpace(Arguments), " ")

	if len(arguments) == 0 || len(arguments[0]) == 0 {
		return "Вы должны передать параметром наименование проекта", nil
	}

	projects, _, err := git.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search:           gitlab.String(arguments[0]),
		SearchNamespaces: gitlab.Bool(true),
	})

	if err != nil {
		log.Print(err)
		return "Произошла непредвиденная ошибка!", nil
	}

	if len(projects) == 0 {
		return "Не было найдено ни единого проекта", nil
	}

	return fmt.Sprintf("TODO: Подписка на проект %s", projects[0].Name), projects[0]
}
