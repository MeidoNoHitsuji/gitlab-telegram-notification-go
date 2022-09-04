package callbacks

type CallbackFuncName string

const (
	TomatoFailNameFuncName CallbackFuncName = "tomato_fail"
	SubscribesFuncName     CallbackFuncName = "subscribes"
	//SelectProjectFuncName  CallbackFuncName = "select_project"
	BackFuncName CallbackFuncName = "back"
)

type DefaultType struct {
	FuncName CallbackFuncName `json:"func_name"`
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

//type SelectProjectType struct {
//	DefaultType
//}
//
//func NewSelectProjectType() *SelectProjectType {
//	return &SelectProjectType{
//		DefaultType: DefaultType{
//			FuncName: SelectProjectFuncName,
//		},
//	}
//}

type BackType struct {
	DefaultType
}

func NewBackType() *BackType {
	return &BackType{
		DefaultType: DefaultType{
			FuncName: BackFuncName,
		},
	}
}
