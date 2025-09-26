package domain

type MockService interface {
	Create(mock *Mock) (*Mock, error)
	Read(id int64) (*Mock, error)
	Update(mock *Mock) (*Mock, error)
	Delete(id int64) error
	List(query *MockQuery) (*MockList, error)
}
