package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab-telegram-notification-go/gitclient"
	"html/template"
	"io/ioutil"
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

	w.Write([]byte(strings.ReplaceAll(strings.TrimSpace(data.Body()), "\n", "<br>")))
	w.WriteHeader(200)
}

func GetWebToggle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["user_id"] == "" {
		http.Error(w, "Вы не передали параметр user_id", http.StatusUnprocessableEntity)
		return
	}

	w.Write([]byte("Hello!"))

	w.WriteHeader(200)
}

func WebToggle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["user_id"] == "" {
		http.Error(w, "Вы не передали параметр user_id", http.StatusUnprocessableEntity)
		return
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("ReadAll")
		fmt.Println(err)
		w.WriteHeader(403)
		return
	}

	var result ToggleData

	err = json.Unmarshal(body, &result)

	if err != nil {
		fmt.Println("Unmarshal")
		fmt.Println(body)
		fmt.Println(err)
		w.WriteHeader(403)
		return
	}

	if result.Payload == "ping" {
		type Response struct {
			ValCode string `json:"validation_code"`
		}

		response := Response{
			ValCode: result.ValidationCode,
		}

		res, err := json.Marshal(response)

		if err != nil {
			fmt.Println("Marshal")
			fmt.Println(err)
			w.WriteHeader(403)
			return
		}

		w.Write(res)
	}

	w.WriteHeader(200)
}
