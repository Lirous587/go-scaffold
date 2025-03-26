package user

import "scaffold/api/mock2"

type ControllerV1 struct {
}

func NewV1() mock2.IMock2V1 {
	controller := &ControllerV1{}
	return controller
}
