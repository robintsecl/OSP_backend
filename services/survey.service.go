package services

import (
	"github.com/robintsecl/osp_backend/models"
)

type SurveyService interface {
	CreateSurvey(*models.Survey) (*string, error)
	GetSurvey(*string) (*models.Survey, error)
	GetAll() ([]*models.Survey, error)
	UpdateSurvey(*models.Survey) error
	DeleteSurvey(*string) error
}
