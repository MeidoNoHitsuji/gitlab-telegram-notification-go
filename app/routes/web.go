package routes

import (
	"github.com/gorilla/mux"
	"gitlab-telegram-notification-go/gitclient"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func WebIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("static/index.html")
	err := t.Execute(w, nil)
	if err != nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func WebPipeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["project_id"] == "" {
		http.Error(w, "Вы не передали параметр project_id", http.StatusUnprocessableEntity)
		return
	}

	if vars["pipeline_id"] == "" {
		http.Error(w, "Вы не передали параметр pipeline_id", http.StatusUnprocessableEntity)
		return
	}

	git := gitclient.Instant()

	projectId, _ := strconv.Atoi(vars["project_id"])
	pipelineId, _ := strconv.Atoi(vars["pipeline_id"])

	pipeline, _, err := git.Pipelines.GetPipeline(projectId, pipelineId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	commits, err := gitclient.GetCommitsLastPipeline(projectId, pipeline.BeforeSHA, pipeline.SHA)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := gitclient.PipelineLogType{
		Commits: commits,
	}

	w.Write([]byte(strings.ReplaceAll(data.Body(), "\n", "</br>")))
	w.WriteHeader(200)
}
