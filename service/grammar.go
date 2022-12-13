package service

import (
	"Bananza/db"
	"Bananza/models"
	"context"
	"time"
)

type GrammarService struct {
	repo db.Grammar
}

func NewGrammarService(repo db.Grammar) *GrammarService {
	return &GrammarService{repo: repo}
}

func (s *GrammarService) GetGrammar(inputGrammar models.InputDictionary) (*[]models.Grammar, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	grammar, err := s.repo.GetGrammar(ctx, inputGrammar)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return grammar, nil
}

func (s *GrammarService) GetDictionary(inputDictionary models.InputDictionary) (*[]models.Dictionary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	dictionary, err := s.repo.GetDictionary(ctx, inputDictionary)
	if err != nil {
		return nil, err
	}
	defer cancel()

	return dictionary, nil
}
