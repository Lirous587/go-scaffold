package mock1

import "scaffold/api/mock1"

type ControllerV1 struct {
}

func NewV1() mock1.IMock1V1 {
	controller := &ControllerV1{}
	return controller
}
