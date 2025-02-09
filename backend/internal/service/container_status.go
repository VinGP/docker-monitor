package service

import (
	"backend/internal/model"
	"backend/internal/repo"
)

type ContainerStatusService struct {
	r *repo.ContainerStatusRepo
}

func NewContainerStatusService(r *repo.ContainerStatusRepo) *ContainerStatusService {
	return &ContainerStatusService{r}
}

func (s *ContainerStatusService) UpsertContainerStatus(status *model.ContainerStatus) error {
	return s.r.UpsertContainerStatus(status)
}

func (s *ContainerStatusService) GetAll() ([]model.ContainerStatus, error) {
	return s.r.GetAll()
}

func (s *ContainerStatusService) DeleteAll() error {
	return s.r.DeleteAll()
}
