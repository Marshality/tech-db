package service

type Repository interface {
	Status() (uint64, uint64, uint64, uint64, error)
	Clear() error
}
