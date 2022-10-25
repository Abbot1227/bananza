package service

type Authorization interface {
}

type Progress interface {
}

type Service struct {
	Authorization
	Progress
}

func NewService() *Service {
	return &Service{}
}
