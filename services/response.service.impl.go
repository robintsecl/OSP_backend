package services

import (
	"context"

	customErr "github.com/robintsecl/osp_backend/errors"
	"github.com/robintsecl/osp_backend/models"

	utils "github.com/robintsecl/osp_backend/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResponseServiceImpl struct {
	responsecollection *mongo.Collection
	surveycollection   *mongo.Collection
	ctx                context.Context
}

func NewResponseService(responsecollection *mongo.Collection, surveycollection *mongo.Collection, ctx context.Context) ResponseService {
	return &ResponseServiceImpl{
		responsecollection: responsecollection,
		surveycollection:   surveycollection,
		ctx:                ctx,
	}
}

func (rsi *ResponseServiceImpl) CreateResponse(response *models.Response) error {
	var survey *models.Survey
	query := bson.D{bson.E{Key: "token", Value: response.Token}}
	err := rsi.surveycollection.FindOne(rsi.ctx, query).Decode(&survey)
	if err != nil {
		return err
	}
	inputErr := utils.ResponseInputChecking(survey.Questions, response.ResponseAnswer)
	if inputErr != nil {
		return inputErr
	}
	_, insertErr := rsi.responsecollection.InsertOne(rsi.ctx, response)
	if insertErr != nil {
		return insertErr
	}
	return nil
}

func (rsi *ResponseServiceImpl) GetAll() ([]*models.Response, error) {
	var responses []*models.Response
	cursor, err := rsi.responsecollection.Find(rsi.ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cursor.Next(rsi.ctx) {
		var response models.Response
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(rsi.ctx)

	return responses, nil
}

func (rsi *ResponseServiceImpl) GetByToken(token *string) ([]*models.Response, error) {
	var responses []*models.Response
	cursor, err := rsi.responsecollection.Find(rsi.ctx, bson.D{bson.E{Key: "token", Value: token}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(rsi.ctx) {
		var response models.Response
		err := cursor.Decode(&response)
		if err != nil {
			return nil, err
		}
		responses = append(responses, &response)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(rsi.ctx)

	return responses, nil
}

func (rsi *ResponseServiceImpl) BatchDeleteByToken(token *string) error {
	query := bson.D{bson.E{Key: "token", Value: token}}
	result, _ := rsi.responsecollection.DeleteMany(rsi.ctx, query)
	if result.DeletedCount < 1 {
		return customErr.ErrDataNotFound
	}
	return nil
}
