package callbacks

type CallbackFuncName string

const (
	StartFuncName                 CallbackFuncName = "start"
	SubscribesFuncName            CallbackFuncName = "subscribes"
	BackFuncName                  CallbackFuncName = "back"
	TomatoFailNameFuncName        CallbackFuncName = "tf_func"   //tomato_fail
	UserSettingsFuncName          CallbackFuncName = "us_func"   //user_settings
	UserSettingTokensFuncName     CallbackFuncName = "ust_func"  //user_setting_tokens
	UserIntegrationsFuncName      CallbackFuncName = "ui_func"   //user_integrations
	UserSettingEnterTokenFuncName CallbackFuncName = "uset_func" //user_setting_enter_token
	SelectProjectSettingsFuncName CallbackFuncName = "sps_func"  //select_project_settings
	SelectFilterFuncName          CallbackFuncName = "sf_func"   //select_filter
	ChangeActiveFuncName          CallbackFuncName = "ca_func"   //change_active
	CreateFilterFuncName          CallbackFuncName = "cf_func"   //create_filter
	EditFilterFuncName            CallbackFuncName = "ef_func"   //edit_filter
	EditFilterParameterFuncName   CallbackFuncName = "efp_func"  //edit_filter_parameter
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

type StartType struct {
	DefaultType
}

func NewStartType() *StartType {
	return &StartType{
		DefaultType: DefaultType{
			FuncName: StartFuncName,
		},
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

type UserSettingsType struct {
	DefaultType
}

func NewUserSettingsType() *UserSettingsType {
	return &UserSettingsType{
		DefaultType: DefaultType{
			FuncName: UserSettingsFuncName,
		},
	}
}

type UserIntegrationsType struct {
	DefaultType
	IntegrationType string `json:"it"`
}

func NewUserIntegrationsType() *UserIntegrationsType {
	return &UserIntegrationsType{
		DefaultType: DefaultType{
			FuncName: UserIntegrationsFuncName,
		},
	}
}

func NewUserIntegrationsWithTypeType(typeIntegration string) *UserIntegrationsType {
	return &UserIntegrationsType{
		DefaultType: DefaultType{
			FuncName: UserIntegrationsFuncName,
		},
		IntegrationType: typeIntegration,
	}
}

type UserSettingTokensType struct {
	DefaultType
	TokenType  string `json:"tt"`
	TokenValue string `json:"tv"`
}

func NewUserSettingTokensType() *UserSettingTokensType {
	return &UserSettingTokensType{
		DefaultType: DefaultType{
			FuncName: UserSettingTokensFuncName,
		},
	}
}

func NewUserSettingTokensWithTokenType(tokenType string) *UserSettingTokensType {
	return &UserSettingTokensType{
		DefaultType: DefaultType{
			FuncName: UserSettingTokensFuncName,
		},
		TokenType: tokenType,
	}
}

type UserSettingEnterTokenType struct {
	DefaultType
	TokenType string `json:"tt"`
}

func NewUserSettingEnterTokenType(tokenType string) *UserSettingEnterTokenType {
	return &UserSettingEnterTokenType{
		DefaultType: DefaultType{
			FuncName: UserSettingEnterTokenFuncName,
		},
		TokenType: tokenType,
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

type EditFilterType struct {
	DefaultType
	ProjectId      int    `json:"pi"`
	EventId        uint   `json:"ei"`
	EventName      string `json:"en"`
	ParameterName  string `json:"pn"`
	ParameterValue string `json:"pv"`
	DeleteValue    bool   `json:"dv"`
	EditFormatter  bool   `json:"ef"`
	FormatterValue string `json:"fv"`
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

type EditFilterWithFormatterType struct {
	EditFilterByIdType
	EditFormatter bool `json:"ef"`
}

type EditFilterWithFormatterValueType struct {
	EditFilterWithFormatterType
	FormatterValue string `json:"fv"`
}

type EditFilterWithDeleteFormatterType struct {
	EditFilterWithFormatterType
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

func NewEditFilterWithFormatterType(projectId int, eventId uint) *EditFilterWithFormatterType {
	return &EditFilterWithFormatterType{
		EditFilterByIdType: EditFilterByIdType{
			DefaultType: DefaultType{
				FuncName: EditFilterFuncName,
			},
			ProjectId: projectId,
			EventId:   eventId,
		},
		EditFormatter: true,
	}
}

func NewEditFilterWithFormatterValueType(projectId int, eventId uint, formatterValue string) *EditFilterWithFormatterValueType {
	return &EditFilterWithFormatterValueType{
		EditFilterWithFormatterType: EditFilterWithFormatterType{
			EditFilterByIdType: EditFilterByIdType{
				DefaultType: DefaultType{
					FuncName: EditFilterFuncName,
				},
				ProjectId: projectId,
				EventId:   eventId,
			},
			EditFormatter: true,
		},
		FormatterValue: formatterValue,
	}
}

func NewEditFilterWithDeleteFormatterType(projectId int, eventId uint) *EditFilterWithDeleteFormatterType {
	return &EditFilterWithDeleteFormatterType{
		EditFilterWithFormatterType: EditFilterWithFormatterType{
			EditFilterByIdType: EditFilterByIdType{
				DefaultType: DefaultType{
					FuncName: EditFilterFuncName,
				},
				ProjectId: projectId,
				EventId:   eventId,
			},
			EditFormatter: true,
		},
		DeleteValue: true,
	}
}

//NewEditFilterWithParameterType (FuncName = EditFilterFuncName)
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

// NewEditFilterParameterType (FuncName = EditFilterParameterFuncName)
func NewEditFilterParameterType(projectId int, eventId uint, parameterName string) *EditFilterWithParameterType {
	return &EditFilterWithParameterType{
		EditFilterByIdType: EditFilterByIdType{
			DefaultType: DefaultType{
				FuncName: EditFilterParameterFuncName,
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
