package callbacks

type CallbackFuncName string

const (
	TomatoFailNameFuncName        CallbackFuncName = "tf_func" //tomato_fail
	SubscribesFuncName            CallbackFuncName = "subscribes"
	BackFuncName                  CallbackFuncName = "back"
	SelectProjectSettingsFuncName CallbackFuncName = "sps_func"       //select_project_settings
	SelectFilterFuncName          CallbackFuncName = "sf_func"        //select_filter
	ChangeActiveFuncName          CallbackFuncName = "ca_func"        //change_active
	CreateFilterFuncName          CallbackFuncName = "cf_func"        //create_filter
	ChoiceWebhookFilterFuncName   CallbackFuncName = "webhook_filter" //TODO: Создать callback На это
	EditFilterFuncName            CallbackFuncName = "ef_func"        //edit_filter
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
	ProjectId     int  `json:"pi"`
	DeleteEventId uint `json:"dei"`
}

func NewSelectProjectSettingsType(projectId int) *SelectProjectSettingsType {
	return &SelectProjectSettingsType{
		DefaultType: DefaultType{
			FuncName: SelectProjectSettingsFuncName,
		},
		ProjectId: projectId,
	}
}

func NewSelectProjectSettingsWithDeleteEventType(projectId int, eventId uint) *SelectProjectSettingsType {
	return &SelectProjectSettingsType{
		DefaultType: DefaultType{
			FuncName: SelectProjectSettingsFuncName,
		},
		ProjectId:     projectId,
		DeleteEventId: eventId,
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
	ProjectId      int    `json:"pi"`
	EventId        uint   `json:"ei"`
	EventName      string `json:"en"`
	ParameterName  string `json:"pn"`
	ParameterValue string `json:"pv"`
	DeleteValue    bool   `json:"dv"`
}

type EditFilterByNameType struct {
	DefaultType
	ProjectId int    `json:"pi"`
	EventName string `json:"en"`
}

type EditFilterByIdType struct {
	DefaultType
	ProjectId int  `json:"pi"`
	EventId   uint `json:"ei"`
}

type EditFilterWithParameterType struct {
	EditFilterByIdType
	ParameterName string `json:"pn"`
}

type EditFilterWithParameterValueType struct {
	EditFilterWithParameterType
	ParameterValue string `json:"pv"`
}

type EditFilterWithDeleteValueType struct {
	EditFilterWithParameterType
	DeleteValue bool `json:"dv"`
}

func NewEditFilterWithEventIdType(projectId int, eventId uint) *EditFilterByIdType {
	return &EditFilterByIdType{
		DefaultType: DefaultType{
			FuncName: EditFilterFuncName,
		},
		ProjectId: projectId,
		EventId:   eventId,
	}
}

func NewEditFilterWithEventNameType(projectId int, eventName string) *EditFilterByNameType {
	return &EditFilterByNameType{
		DefaultType: DefaultType{
			FuncName: EditFilterFuncName,
		},
		ProjectId: projectId,
		EventName: eventName,
	}
}

func NewEditFilterWithParameterType(projectId int, eventId uint, parameterName string) *EditFilterWithParameterType {
	return &EditFilterWithParameterType{
		EditFilterByIdType: EditFilterByIdType{
			DefaultType: DefaultType{
				FuncName: EditFilterFuncName,
			},
			ProjectId: projectId,
			EventId:   eventId,
		},
		ParameterName: parameterName,
	}
}

func NewEditFilterWithParameterValueType(projectId int, eventId uint, parameterName string, parameterValue string) *EditFilterWithParameterValueType {
	return &EditFilterWithParameterValueType{
		EditFilterWithParameterType: EditFilterWithParameterType{
			EditFilterByIdType: EditFilterByIdType{
				DefaultType: DefaultType{
					FuncName: EditFilterFuncName,
				},
				ProjectId: projectId,
				EventId:   eventId,
			},
			ParameterName: parameterName,
		},
		ParameterValue: parameterValue,
	}
}

func NewEditFilterWithDeleteParameterType(projectId int, eventId uint, parameterName string) *EditFilterWithDeleteValueType {
	return &EditFilterWithDeleteValueType{
		EditFilterWithParameterType: EditFilterWithParameterType{
			EditFilterByIdType: EditFilterByIdType{
				DefaultType: DefaultType{
					FuncName: EditFilterFuncName,
				},
				ProjectId: projectId,
				EventId:   eventId,
			},
			ParameterName: parameterName,
		},
		DeleteValue: true,
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
