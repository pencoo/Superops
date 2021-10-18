package factory

import "testing"

//通过调用传递不同参数，以运行不同的工厂，实现不同的功能

func TestNewfactory(t *testing.T) {
	Newfactory("Kentucky").Cook()
	Newfactory("McDonalds").Cook()
}
