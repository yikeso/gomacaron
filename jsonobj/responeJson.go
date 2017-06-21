package jsonobj

//操作状态，和错误信息
type BaseRespone struct {
	Status int //0为操作成功
	Message string
}
//返回一个结构体
type OneRespone struct {
	BaseRespone
	One interface{}
}