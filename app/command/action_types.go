package command

type ActionNameType string

const (
	SelectProject       ActionNameType = "select_project"
	SelectProjectOption ActionNameType = "select_project_option"
	SelectWebHookEvent  ActionNameType = "select_webhook_event"
	SelectFilter        ActionNameType = "select_filter"
	EditParameterFilter ActionNameType = "edit_parameter_filter"
)

type ActionInitByType string

const (
	InitByCommand ActionInitByType = "command"
	InitByText    ActionInitByType = "text"
	InitByAction  ActionInitByType = "action"
)

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
	// Доступные варианты: InitByCommand, InitByText, InitByAction
	InitBy ActionInitByType

	// InitText это текст, которым тригерится данный экшн при помощи типа, который указан в InitBy.
	InitText string

	// AfterAction это условие, которое оозначающее предыдущий Action
	//
	// Пустое значение - должно отсутствовать
	// nil - любой action
	// ActionNameType - соответствует имени экшна
	AfterAction string
}

// Active это действие, которое выполняется, когда был получен Callback или Message и используется текущий BaseAction.
func (act *BaseAction) Active(value string) error {
	//TODO
	return nil
}

// SelectProjectAction вызывается когда надо выбрать проект для редактирования
type SelectProjectAction struct {
	ID           ActionNameType
	BeforeAction ActionNameType
	InitBy       ActionInitByType
	InitText     string
	AfterAction  string
}

func (act *SelectProjectAction) New() *SelectProjectAction {
	act.ID = SelectProject
	act.InitBy = InitByCommand
	act.InitText = "subscribe"
	act.AfterAction = ""
	return act
}

func (act *SelectProjectAction) Active(value string) error {
	//TODO
	return nil
}

// SelectProjectOptionAction вызывается когда надо выбирать действие для проекта
type SelectProjectOptionAction struct {
	ID           ActionNameType
	BeforeAction ActionNameType
	InitBy       ActionInitByType
	InitText     string
	AfterAction  string
}

func (act SelectProjectOptionAction) New() {
	act.ID = SelectProjectOption
	act.BeforeAction = SelectProject
}

// SelectWebHookEventAction вызывается когда надо выбрать имя ивента
type SelectWebHookEventAction struct {
	ID           ActionNameType
	BeforeAction ActionNameType
	InitBy       ActionInitByType
	InitText     string
	AfterAction  string
}

func (act SelectWebHookEventAction) New() {
	act.ID = SelectWebHookEvent
	act.BeforeAction = SelectProjectOption
}

// SelectFilterActon выбор соответствующего фильтра. Или же создание нового.
type SelectFilterActon struct {
	ID           ActionNameType
	BeforeAction ActionNameType
	InitBy       ActionInitByType
	InitText     string
	AfterAction  string
}

func (act SelectFilterActon) New() {
	act.ID = SelectFilter
	act.BeforeAction = SelectWebHookEvent
}

// EditParameterFilterAction редактирование параметров фильтра
type EditParameterFilterAction struct {
	ID           ActionNameType
	BeforeAction ActionNameType
	InitBy       ActionInitByType
	InitText     string
	AfterAction  string
}

func (act EditParameterFilterAction) New() {
	act.ID = EditParameterFilter
	act.BeforeAction = SelectFilter
}
