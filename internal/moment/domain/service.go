package domain

type MomentService interface {
	Create(moment *Moment) (*Moment, error)
	Delete(id int64) error
	Update(moment *Moment) (*Moment, error)
	ListMoments(query *MomentQuery) (*MomentPages, error)
	RandomN(count int8) ([]*Moment, error)
}
