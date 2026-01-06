package service

type IService interface {
}

type Service struct {
}

func NewService() IService {
	return &Service{}
}

func (s *Service) Func1() string {
	return "Service Func1"
}

func (s *Service) Func2() string {
	return "Service Func2"
}
