package telegram

const (
	TomatoFailName = "tomato_fail"
)

type DefaultType struct {
	FuncName string `json:"func_name"`
}

type TomatoFailType struct {
	DefaultType
	Count int `json:"count"`
}

func NewTomatoFailType(count int) TomatoFailType {
	return TomatoFailType{
		DefaultType: DefaultType{
			FuncName: TomatoFailName,
		},
		Count: count,
	}
}
