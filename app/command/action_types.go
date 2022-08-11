package command

// BaseAction это абстрактный тип для всех экшнов.
type BaseAction struct {
	//ID это уникальный идентификатор экшна. Оно записывается в базу.
	ID string

	//InitBy это условие по которому триггерится данный экшн.
	//Доступные варианты: command, text, action
	InitBy string

	//InitText это текст, которым тригерится данный экшн при помощи типа, который указан в InitBy.
	InitText string
}

// Action это действие, которое выполняется, когда был получен Callback или Message и используется текущий BaseAction.
func (act *BaseAction) Action(value string) error {
	//TODO
	return nil
}
