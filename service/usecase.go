package service

type Usecase interface {
	GetStatus() (uint64, uint64, uint64, uint64, error)
	DoClear() error
}
