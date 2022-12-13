package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab-telegram-notification-go/gitclient"
	"gitlab-telegram-notification-go/jiraclient"
	"gitlab-telegram-notification-go/routes/request"
	"gitlab-telegram-notification-go/toggl"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	limited = make(map[int64]time.Time)
)

func WebIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Hello World!!",
	})
}

func WebPipeline(c *gin.Context) {
	if c.Param("project_id") == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Вы не передали параметр project_id",
		})
		return
	}

	if c.Param("pipeline_id") == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Вы не передали параметр pipeline_id",
		})
		return
	}

	git := gitclient.Instant()

	projectId, _ := strconv.Atoi(c.Param("project_id"))
	pipelineId, _ := strconv.Atoi(c.Param("pipeline_id"))

	pipeline, _, err := git.Pipelines.GetPipeline(projectId, pipelineId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	beforePipeline, err := gitclient.GetBeforeFailedPipeline(projectId, pipeline.BeforeSHA, pipeline.Ref)
	beforeSHA := pipeline.BeforeSHA

	if err == nil && beforePipeline != nil {
		beforeSHA = beforePipeline.BeforeSHA
	}

	commits, err := gitclient.GetCommitsLastPipeline(projectId, beforeSHA, pipeline.SHA)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	data := gitclient.PipelineLogType{
		Commits: commits,
	}
	c.Data(
		http.StatusOK,
		"text/html; charset=utf-8",
		[]byte(strings.ReplaceAll(strings.TrimSpace(data.Body()), "\n", "<br>")),
	)
}

func WebToggle(c *gin.Context) {
	signature := c.GetHeader("X-Webhook-Signature-256")

	if signature == "" {
		c.String(http.StatusBadRequest, "X-Webhook-Signature-256 not found")
		return
	}

	body, err := c.GetRawData()

	if err != nil {
		fmt.Println("Body not found")
		c.String(
			http.StatusBadRequest,
			"Body not found",
		)
		return
	}

	secret := os.Getenv("TOGGLE_SECRET")

	fmt.Println(string(body))

	if !toggl.HmacIsValid(string(body), signature, secret) {
		fmt.Println("Unauthorized")
		c.String(
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	if c.Param("user_telegram_id") == "" {
		c.String(
			http.StatusUnprocessableEntity,
			"user_telegram_id not found",
		)
		return
	}

	telegramChannelId, err := strconv.ParseInt(c.Param("user_telegram_id"), 10, 64)

	if err != nil {
		c.String(
			http.StatusUnprocessableEntity,
			"user_telegram_id not found",
		)
		return
	}

	var result request.ValidationData

	err = json.Unmarshal(body, &result)

	if err == nil && result.Payload == "ping" {
		type Response struct {
			ValCode string `json:"validation_code"`
		}

		response := Response{
			ValCode: result.ValidationCode,
		}

		res, err := json.Marshal(response)

		if err != nil {
			c.String(
				http.StatusBadRequest,
				fmt.Sprintf("Bad Body: %s", err.Error()),
			)
			return
		}

		c.Data(
			http.StatusOK,
			fmt.Sprintf("%s; charset=utf-8", binding.MIMEJSON),
			res,
		)
		return
	} else if err == nil {
		c.Status(http.StatusOK)
		fmt.Println("Неизвестный запрос!!")
		return
	}

	var data request.ToggleData

	err = json.Unmarshal(body, &data)

	if err != nil {
		c.String(
			http.StatusBadRequest,
			fmt.Sprintf("Bad Body: %s", err.Error()),
		)
		return
	}

	if data.Metadata.Action == "updated" {
		jiraclient.UpdateJiraWorklog(telegramChannelId, data)
	} else if data.Metadata.Action == "deleted" {
		jiraclient.DeleteJiraWorklog(telegramChannelId, data)
	}

	c.Status(http.StatusOK)
}
