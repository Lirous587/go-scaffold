package domain

type MomentCoordinate struct {
	X float64 // 经度 (longitude)
	Y float64 // 纬度 (latitude)
}

const InvalidCoordinate = -999.0

func (coor MomentCoordinate) IsValid() bool {
	return coor.X != InvalidCoordinate && coor.Y != InvalidCoordinate
}

func NewInvalidMomentCoordinate() *MomentCoordinate {
	return &MomentCoordinate{X: InvalidCoordinate, Y: InvalidCoordinate}
}
