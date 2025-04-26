package services

import "github.com/robintsecl/osp_backend/models"

type ResponseService interface {
	CreateResponse(response *models.Response) error
	GetAll() ([]*models.Response, error)
	GetByToken(token *string) ([]*models.Response, error)
	BatchDeleteByToken(token *string) error
}
