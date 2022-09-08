package callbacks

type CallbackFuncName string

const (
	TomatoFailNameFuncName        CallbackFuncName = "tomato_fail"
	SubscribesFuncName            CallbackFuncName = "subscribes"
	BackFuncName                  CallbackFuncName = "back"
	SelectProjectSettingsFuncName CallbackFuncName = "select_project_settings"
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
