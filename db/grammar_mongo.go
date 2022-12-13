package db

import (
	"Bananza/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type GrammarMongo struct{}

func NewGrammarMongo() *GrammarMongo {
	return &GrammarMongo{}
}

func (r *GrammarMongo) GetDictionary(ctx context.Context, inputDictionary models.InputDictionary) (*[]models.Dictionary, error) {
	var dictionary []models.Dictionary
	filter := bson.D{
		{"$and",
			bson.A{
				bson.D{{"level", inputDictionary.Level}},
				bson.D{{"language", inputDictionary.Language}},
			},
		},
	}

	cursor, err := dictionaryCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &dictionary); err != nil {
		return nil, err
	}

	return &dictionary, nil
}
