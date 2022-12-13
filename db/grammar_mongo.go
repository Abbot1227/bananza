package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type GrammarMongo struct{}

func NewGrammarMongo() *GrammarMongo {
	return &GrammarMongo{}
}

func (r *GrammarMongo) GetGrammar(ctx context.Context, inputGrammar models.InputDictionary) (*[]models.Grammar, error) {
	var grammar []models.Grammar
	matchLanguageStage := bson.D{{"$match", bson.D{{"language", inputGrammar.Language}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", inputGrammar.Level}}}}}}

	cursor, err := grammarCollection.Find(ctx, mongo.Pipeline{matchLanguageStage, matchLevelStage})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &grammar); err != nil {
		return nil, err
	}
	return &grammar, nil
}

func (r *GrammarMongo) GetDictionary(ctx context.Context, inputDictionary models.InputDictionary) (*[]models.Dictionary, error) {
	var dictionary []models.Dictionary
	matchLanguageStage := bson.D{{"$match", bson.D{{"language", inputDictionary.Language}}}}
	matchLevelStage := bson.D{{"$match", bson.D{{"level", bson.D{{"$lte", inputDictionary.Level}}}}}}

	cursor, err := dictionaryCollection.Find(ctx, mongo.Pipeline{matchLanguageStage, matchLevelStage})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &dictionary); err != nil {
		return nil, err
	}

	return &dictionary, nil
}
