package abstract_factory

import "fmt"

type Api interface {
	Cook()
}

type Kentucky struct{}

func (*Kentucky) Cook() {
	fmt.Println("cooking Kentucky...")
}

type McDonalds struct{}

func (*McDonalds) Cook() {
	fmt.Println("cooking McDonalds...")
}

//抽象工厂接口
type AbstractFactory interface {
	CreateKentucky() Api
	CreateMcDonalds() Api
}

//实现工厂接口
type SimpleCookAPI struct{}

func NewSimpleCookAPI() AbstractFactory {
	return &SimpleCookAPI{}
}

func (*SimpleCookAPI) CreateKentucky() Api {
	return &Kentucky{}
}
func (*SimpleCookAPI) CreateMcDonalds() Api {
	return &McDonalds{}
}
