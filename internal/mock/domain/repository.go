package domain


type MockRepository interface {
	FindByID(id int64) (*Mock, error)

	Create(mock *Mock) (*Mock, error)
	Update(mock *Mock) (*Mock, error)
	Delete(id int64) error
	List(query *MockQuery) (*MockList, error)
}

type MockCache interface {

}
