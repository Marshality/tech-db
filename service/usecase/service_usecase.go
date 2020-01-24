package usecase

import (
	"github.com/Marshality/tech-db/service"
)

type ServiceUsecase struct {
	serviceRepo service.Repository
}

func NewServiceUsecase(sr service.Repository) service.Usecase {
	return &ServiceUsecase{
		serviceRepo: sr,
	}
}

func (su *ServiceUsecase) GetStatus() (uint64, uint64, uint64, uint64, error) {
	return su.serviceRepo.Status()
}

func (su *ServiceUsecase) DoClear() error {
	return su.serviceRepo.Clear()
}
