package services

import (
	"context"
	"math/rand"

	"github.com/robintsecl/osp_backend/models"
	utils "github.com/robintsecl/osp_backend/utils"

	customErr "github.com/robintsecl/osp_backend/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: remove usercollection

type SurveyServiceImpl struct {
	surveycollection *mongo.Collection
	ctx              context.Context
}

func NewSurveyService(surveycollection *mongo.Collection, ctx context.Context) SurveyService {
	return &SurveyServiceImpl{
		surveycollection: surveycollection,
		ctx:              ctx,
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func GenerateRandomToken() string {
	tokenRunes := make([]rune, 5)
	for i := range tokenRunes {
		tokenRunes[i] = letters[rand.Intn(len(letters))]
	}
	return string(tokenRunes)
}

func (ssi *SurveyServiceImpl) CreateSurvey(survey *models.Survey) (*string, error) {
	// Check if token alreay exists
	for range 20 {
		var existingToken *string
		var tokenToCheck = GenerateRandomToken()
		query := bson.D{bson.E{Key: "token", Value: tokenToCheck}}
		err := ssi.surveycollection.FindOne(ssi.ctx, query).Decode(&existingToken)
		if err != nil && err != mongo.ErrNoDocuments {
			return nil, customErr.ErrDBBadGateway
		}
		if existingToken == nil {
			survey.Token = tokenToCheck
			break
		}
	}

	commonErr := utils.CommonChecking(&survey.Questions)
	if commonErr != nil {
		return nil, commonErr
	}
	_, err := ssi.surveycollection.InsertOne(ssi.ctx, survey)
	if err != nil {
		return nil, err
	}
	return &survey.Token, nil
}

func (ssi *SurveyServiceImpl) GetSurvey(token *string) (*models.Survey, error) {
	var survey *models.Survey
	query := bson.D{bson.E{Key: "token", Value: token}}
	err := ssi.surveycollection.FindOne(ssi.ctx, query).Decode(&survey)
	return survey, err
}

func (ssi *SurveyServiceImpl) GetAll() ([]*models.Survey, error) {
	var surveys []*models.Survey
	cursor, err := ssi.surveycollection.Find(ssi.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(ssi.ctx) {
		var survey models.Survey
		err := cursor.Decode(&survey)
		if err != nil {
			return nil, err
		}
		surveys = append(surveys, &survey)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(ssi.ctx)

	return surveys, nil
}

func (ssi *SurveyServiceImpl) UpdateSurvey(survey *models.Survey) error {
	commonErr := utils.CommonChecking(&survey.Questions)
	if commonErr != nil {
		return commonErr
	}
	query := bson.D{bson.E{Key: "token", Value: survey.Token}}
	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "title", Value: survey.Title},
		bson.E{Key: "questions", Value: survey.Questions},
	}}}
	result, _ := ssi.surveycollection.UpdateOne(ssi.ctx, query, update)
	if result.MatchedCount != 1 {
		return customErr.ErrDataNotFound
	}
	return nil
}

func (ssi *SurveyServiceImpl) DeleteSurvey(token *string) error {
	query := bson.D{bson.E{Key: "token", Value: token}}
	result, _ := ssi.surveycollection.DeleteOne(ssi.ctx, query)
	if result.DeletedCount != 1 {
		return customErr.ErrDataNotFound
	}
	return nil
}
