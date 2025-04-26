package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/robintsecl/osp_backend/constants"
	customErr "github.com/robintsecl/osp_backend/errors"
	"github.com/robintsecl/osp_backend/models"
	"github.com/robintsecl/osp_backend/services"
	utils "github.com/robintsecl/osp_backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type SurveyController struct {
	SurveyService  services.SurveyService
	usercollection *mongo.Collection
	validate       *validator.Validate
}

func NewSurveyController(surveyservice services.SurveyService, usercollection *mongo.Collection, validate *validator.Validate) SurveyController {
	return SurveyController{
		SurveyService:  surveyservice,
		usercollection: usercollection,
		validate:       validate,
	}
}

func (sc *SurveyController) CreateSurvey(ctx *gin.Context) {
	// Check admin
	admErr := utils.CheckAdmin(ctx, sc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	// Bind json
	var survey models.Survey
	if err := ctx.ShouldBindJSON((&survey)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// TODO: pass validator from main
	if err := sc.validate.Struct(survey); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Insert date
	utils.InsertSurveyDate(&survey, constants.ACTION_DATE_CREATE)
	// Create
	token, err := sc.SurveyService.CreateSurvey(&survey)
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success with token: " + *token})
}

func (sc *SurveyController) GetSurvey(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		customErr.ThrowCustomError(&customErr.ErrQueryParamMissing, ctx)
		return
	}
	survey, err := sc.SurveyService.GetSurvey(&token)
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, survey)
}

func (sc *SurveyController) GetAll(ctx *gin.Context) {
	// Check admin
	admErr := utils.CheckAdmin(ctx, sc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	surveys, err := sc.SurveyService.GetAll()
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, surveys)
}

func (sc *SurveyController) UpdateSurvey(ctx *gin.Context) {
	admErr := utils.CheckAdmin(ctx, sc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}

	var survey models.Survey
	if err := ctx.ShouldBindJSON((&survey)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := sc.validate.Struct(survey); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	utils.InsertSurveyDate(&survey, constants.ACTION_DATE_UPDATE)
	err := sc.SurveyService.UpdateSurvey(&survey)
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, survey)
}

func (sc *SurveyController) DeleteSurvey(ctx *gin.Context) {
	admErr := utils.CheckAdmin(ctx, sc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	token := ctx.Query("token")
	if token == "" {
		customErr.ThrowCustomError(&customErr.ErrQueryParamMissing, ctx)
		return
	}
	err := sc.SurveyService.DeleteSurvey(&token)
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (sc *SurveyController) RegisterSurveyRoutes(rg *gin.RouterGroup) {
	surveyroute := rg.Group("/survey")
	surveyroute.POST("/create", sc.CreateSurvey)
	surveyroute.GET("/get", sc.GetSurvey)
	surveyroute.GET("/getall", sc.GetAll)
	surveyroute.PUT("/update", sc.UpdateSurvey)
	surveyroute.DELETE("/delete", sc.DeleteSurvey)
}
