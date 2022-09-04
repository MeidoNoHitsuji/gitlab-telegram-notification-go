package actions

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"strings"
)

type ActionNameType string

const (
	SelectProjectOption ActionNameType = "select_project_option"
	SelectWebHookEvent  ActionNameType = "select_webhook_event"
	SelectFilter        ActionNameType = "select_filter"
	EditParameterFilter ActionNameType = "edit_parameter_filter"

	TestActionType ActionNameType = "test"
)

type ActionInitByType string

const (
	InitByCommand  ActionInitByType = "actions"
	InitByText     ActionInitByType = "text"
	InitByCallback ActionInitByType = "callback"
)

type BaseInterface interface {
	Validate(update tgbotapi.Update) bool
	Active(update tgbotapi.Update) error
	GetActionName() ActionNameType
	GetBeforeAction() ActionNameType
}

// BaseAction это абстрактный тип для всех экшнов.
type BaseAction struct {
	// ID это уникальный идентификатор экшна. Оно записывается в базу.
	ID ActionNameType

	// BeforeAction это предшестующий текущему экшн
	//
	// nil - отсутствует кнопка возврата
	// ActionNameType - возвращает на соответствующее состояние экшна
	BeforeAction ActionNameType

	// InitBy это условие по которому триггерится данный экшн.
	// Доступные варианты: InitByCommand, InitByText, InitByCallback
	InitBy ActionInitByType

	// InitText это текст, которым тригерится данный экшн при помощи типа, который указан в InitBy.
	InitText string

	// InitCallbackFuncName это имя функции, которой произошла активация
	InitCallbackFuncName callbacks.CallbackFuncName

	// CallbackData это данные, которые были получены, если action триггернулся на  InitByCallback
	CallbackData *callbacks.DefaultType

	// AfterAction это условие, которое означающее предыдущий Action
	//
	// Пустое значение - должно отсутствовать
	// nil - любой action
	// ActionNameType - соответствует имени экшна
	AfterAction ActionNameType
}

func (act *BaseAction) GetActionName() ActionNameType {
	return act.ID
}

func (act *BaseAction) GetBeforeAction() ActionNameType {
	return act.BeforeAction
}

func (act *BaseAction) Validate(update tgbotapi.Update) bool {

	if act.AfterAction != "" {
		if act.AfterAction != GetActualAction(update) {
			return false
		}
	}

	switch act.InitBy {
	case InitByText:
		if update.Message == nil || update.Message.IsCommand() {
			return false
		}

		return act.InitText == "" || strings.ToLower(update.Message.Text) == strings.ToLower(act.InitText)
	case InitByCommand:
		if update.Message == nil || !update.Message.IsCommand() {
			return false
		}

		return update.Message.Command() == act.InitText
	case InitByCallback:
		if update.CallbackQuery == nil || update.CallbackQuery.Data == "" {
			return false
		}

		err := json.Unmarshal([]byte(update.CallbackQuery.Data), &act.CallbackData)
		if err != nil {
			return false
		}

		return act.CallbackData.FuncName == act.InitCallbackFuncName
	}
	return false
}

// SelectProjectOptionAction вызывается когда надо выбирать действие для проекта
type SelectProjectOptionAction struct {
	BaseAction
}

func NewSelectProjectOptionAction() *SelectProjectOptionAction {
	return &SelectProjectOptionAction{
		BaseAction: BaseAction{
			ID:           SelectProjectOption,
			BeforeAction: SelectProjectActionType,
		},
	}
}

// SelectWebHookEventAction вызывается когда надо выбрать имя ивента
type SelectWebHookEventAction struct {
	BaseAction
}

func NewSelectWebHookEventAction() *SelectWebHookEventAction {
	return &SelectWebHookEventAction{
		BaseAction: BaseAction{
			ID:           SelectWebHookEvent,
			BeforeAction: SelectProjectOption,
		},
	}
}

// SelectFilterActon выбор соответствующего фильтра. Или же создание нового.
type SelectFilterActon struct {
	BaseAction
}

func NewSelectFilterActon() *SelectFilterActon {
	return &SelectFilterActon{
		BaseAction: BaseAction{
			ID:           SelectFilter,
			BeforeAction: SelectWebHookEvent,
		},
	}
}

// EditParameterFilterAction редактирование параметров фильтра
type EditParameterFilterAction struct {
	BaseAction
}

func NewEditParameterFilterAction() *EditParameterFilterAction {
	return &EditParameterFilterAction{
		BaseAction: BaseAction{
			ID:           EditParameterFilter,
			BeforeAction: SelectFilter,
		},
	}
}
