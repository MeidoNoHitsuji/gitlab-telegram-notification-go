package actions

import (
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gitlab-telegram-notification-go/actions/callbacks"
	"strings"
)

type ActionNameType string

const (
	SelectWebHookEvent ActionNameType = "swe_act" //select_webhook_event
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
	SetIsBack(update tgbotapi.Update) error
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
	InitBy []ActionInitByType

	// InitializedBy если InitBy создерживаем множествество элементов, то указывает чем было инициализировано
	InitializedBy ActionInitByType

	// InitText это текст, которым тригерится данный экшн при помощи типа, который указан в InitBy.
	InitText string

	// InitCallbackFuncNames это имя функции, которой произошла активация
	InitCallbackFuncNames []callbacks.CallbackFuncName

	//// CallbackData это данные, которые были получены, если action триггернулся на  InitByCallback
	//CallbackData *callbacks.DefaultType `json:"callback_data"`

	// AfterAction это условие, которое означающее предыдущий Action
	//
	// Пустое значение - должно отсутствовать
	// nil - любой action
	// ActionNameType - соответствует имени экшна
	AfterAction ActionNameType

	// IsBack отвечает является ли выполнение текущего Action'а возвратом
	IsBack bool

	// BackData это даные, которые были получены при кнопке Отмена, если IsBack = true
	BackData any `json:"bd"`
}

func (act *BaseAction) GetActionName() ActionNameType {
	return act.ID
}

func (act *BaseAction) SetIsBack(update tgbotapi.Update) error {
	act.IsBack = true
	return nil
}

func (act *BaseAction) GetBeforeAction() ActionNameType {
	return act.BeforeAction
}

func (act *BaseAction) Validate(update tgbotapi.Update) bool {

	if act.AfterAction != "" {
		actualAction := GetActualAction(update)
		if act.ID != actualAction && act.AfterAction != actualAction {
			return false
		}
	}

	for _, byType := range act.InitBy {
		res := false

		switch byType {
		case InitByText:
			if update.Message == nil || update.Message.IsCommand() {
				continue
			}

			res = act.InitText == "" || strings.ToLower(update.Message.Text) == strings.ToLower(act.InitText)
		case InitByCommand:
			if update.Message == nil || !update.Message.IsCommand() {
				continue
			}

			res = update.Message.Command() == act.InitText
		case InitByCallback:
			if update.CallbackQuery == nil || update.CallbackQuery.Data == "" {
				continue
			}

			var callbackData *callbacks.DefaultType

			// Данное преобразование нужно, чтобы вытянуть имя функции
			err := json.Unmarshal([]byte(update.CallbackQuery.Data), &callbackData)
			if err != nil {
				continue
			}

			for _, a := range act.InitCallbackFuncNames {
				if a == callbackData.FuncName {
					res = true
				}
			}
		}

		if res {
			act.InitializedBy = byType
			return res
		}
	}

	return false
}

// SelectWebHookEventAction вызывается когда надо выбрать имя ивента
type SelectWebHookEventAction struct {
	BaseAction
}

func NewSelectWebHookEventAction() *SelectWebHookEventAction {
	return &SelectWebHookEventAction{
		BaseAction: BaseAction{
			ID:           SelectWebHookEvent,
			BeforeAction: SelectProjectSettings,
		},
	}
}

//// EditParameterFilterAction редактирование параметров фильтра
//type EditParameterFilterAction struct {
//	BaseAction
//}
//
//func NewEditParameterFilterAction() *EditParameterFilterAction {
//	return &EditParameterFilterAction{
//		BaseAction: BaseAction{
//			ID:           EditParameterFilter,
//			BeforeAction: SelectFilter,
//		},
//	}
//}
