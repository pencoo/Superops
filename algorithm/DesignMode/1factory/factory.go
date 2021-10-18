package factory

import "fmt"

func Newfactory(name string) Api {
	switch name {
	case "Kentucky", "kentucky":
		return &Kentucky{}
	case "mcdonalds", "McDonalds":
		return &McDonalds{}
	default:
		return nil
	}
}

type Api interface {
	Cook()
}

type Kentucky struct{}

func (*Kentucky) Cook() {
	fmt.Println("cooking Kentucky...")
}

type McDonalds struct{}

func (*McDonalds) Cook() {
	fmt.Println("cooking mcdonalds...")
}
