package db

type Authorization interface {
}

type Progress interface {
}

type Repository struct {
	Authorization
	Progress
}

func NewRepository() *Repository {
	return &Repository{}
}
