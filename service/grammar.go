package service

import "Bananza/db"

type GrammarService struct {
	repo db.Grammar
}

func NewGrammarService(repo db.Grammar) *GrammarService {
	return &GrammarService{repo: repo}
}
