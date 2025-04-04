package services

import repo "app/repos"

type IService interface {
	PingService() interface{}
}
type ServiceInstance struct {
	repo repo.Repo
}

// easy swith of repo definition
func NewServiceInstance(u repo.Repo) *ServiceInstance {

	s := &ServiceInstance{
		repo: u,
	}
	return s
}
func (s *ServiceInstance) PingService() interface{} {
	return s.repo.Ping()
}
