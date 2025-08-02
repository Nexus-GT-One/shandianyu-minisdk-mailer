package constant

type CheckBSideResult string

const (
	// 可以进入
	ENTER = CheckBSideResult("enter")
	// 不能进入
	CLOSE = CheckBSideResult("close")
)
