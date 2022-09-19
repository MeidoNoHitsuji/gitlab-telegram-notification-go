package callbacks

type CallbackFuncName string

const (
	TomatoFailNameFuncName        CallbackFuncName = "tomato_fail"
	SubscribesFuncName            CallbackFuncName = "subscribes"
	BackFuncName                  CallbackFuncName = "back"
	SelectProjectSettingsFuncName CallbackFuncName = "select_project_settings"
	SelectFilterFuncName          CallbackFuncName = "select_filter"
	ChangeActiveFuncName          CallbackFuncName = "change_active"
	CreateFilterFuncName          CallbackFuncName = "create_filter"
	ChoiceWebhookFilterFuncName   CallbackFuncName = "webhook_filter" //TODO: Создать callback На это
	EditFilterFuncName            CallbackFuncName = "edit_filter"
)

type DefaultType struct {
	FuncName CallbackFuncName `json:"fn"`
}

type TomatoFailType struct {
	DefaultType
	Count int `json:"count"`
}

func NewTomatoFailType(count int) TomatoFailType {
	return TomatoFailType{
		DefaultType: DefaultType{
			FuncName: TomatoFailNameFuncName,
		},
		Count: count,
	}
}

type SubscribesType struct {
	DefaultType
}

func NewSubscribesType() *SubscribesType {
	return &SubscribesType{
		DefaultType: DefaultType{
			FuncName: SubscribesFuncName,
		},
	}
}

type SelectProjectSettingsType struct {
	DefaultType
	ProjectId int `json:"pi"`
}

func NewSelectProjectSettingsType(projectId int) *SelectProjectSettingsType {
	return &SelectProjectSettingsType{
		DefaultType: DefaultType{
			FuncName: SelectProjectSettingsFuncName,
		},
		ProjectId: projectId,
	}
}

type ChangeActiveProjectType struct {
	DefaultType
	ProjectId int `json:"pi"`
}

func NewChangeActiveProjectType(projectId int) *SelectProjectSettingsType {
	return &SelectProjectSettingsType{
		DefaultType: DefaultType{
			FuncName: ChangeActiveFuncName,
		},
		ProjectId: projectId,
	}
}

type SelectFilterType struct {
	DefaultType
	ProjectId int `json:"pi"`
}

func NewSelectFilterType(projectId int) *SelectFilterType {
	return &SelectFilterType{
		DefaultType: DefaultType{
			FuncName: SelectFilterFuncName,
		},
		ProjectId: projectId,
	}
}

type CreateFilterType struct {
	DefaultType
	ProjectId int `json:"pi"`
}

func NewCreateFilterType(projectId int) *CreateFilterType {
	return &CreateFilterType{
		DefaultType: DefaultType{
			FuncName: CreateFilterFuncName,
		},
		ProjectId: projectId,
	}
}

type ChoiceWebhookFilterType struct {
	DefaultType
	ProjectId int `json:"pi"`
}

func NewChoiceWebhookFilterType(projectId int) *ChoiceWebhookFilterType {
	return &ChoiceWebhookFilterType{
		DefaultType: DefaultType{
			FuncName: ChoiceWebhookFilterFuncName,
		},
		ProjectId: projectId,
	}
}

type EditFilterType struct {
	DefaultType
	ProjectId int    `json:"pi"`
	EventId   uint   `json:"ei"`
	EventName string `json:"en"`
}

func NewEditFilterWithEventIdType(projectId int, eventId uint) *EditFilterType {
	return &EditFilterType{
		DefaultType: DefaultType{
			FuncName: EditFilterFuncName,
		},
		ProjectId: projectId,
		EventId:   eventId,
	}
}

func NewEditFilterWithEventNameType(projectId int, eventName string) *EditFilterType {
	return &EditFilterType{
		DefaultType: DefaultType{
			FuncName: EditFilterFuncName,
		},
		ProjectId: projectId,
		EventName: eventName,
	}
}

type BackType struct {
	DefaultType
	BackData any `json:"bd"`
}

func NewBackType(backData any) *BackType {
	return &BackType{
		DefaultType: DefaultType{
			FuncName: BackFuncName,
		},
		BackData: backData,
	}
}
