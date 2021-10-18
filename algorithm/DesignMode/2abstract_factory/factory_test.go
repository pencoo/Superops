package abstract_factory

import "testing"

func TestNewSimpleCookAPI(t *testing.T) {
	factory := NewSimpleCookAPI()
	k := factory.CreateKentucky()
	k.Cook()

	m := factory.CreateMcDonalds()
	m.Cook()
}
