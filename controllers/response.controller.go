package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	constants "github.com/robintsecl/osp_backend/constants"
	customErr "github.com/robintsecl/osp_backend/errors"
	"github.com/robintsecl/osp_backend/models"
	"github.com/robintsecl/osp_backend/services"
	utils "github.com/robintsecl/osp_backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResponseController struct {
	ResponseService services.ResponseService
	usercollection  *mongo.Collection
	validate        *validator.Validate
}

func NewResponseController(responseservice services.ResponseService, usercollection *mongo.Collection, validate *validator.Validate) ResponseController {
	return ResponseController{
		ResponseService: responseservice,
		usercollection:  usercollection,
		validate:        validate,
	}
}

func (rc *ResponseController) CreateResponse(ctx *gin.Context) {
	fmt.Printf("Creating response!\n")
	// Bind json
	var response models.Response
	if err := ctx.ShouldBindJSON((&response)); err != nil {
		fmt.Printf("Failed to bind JSON!\n")
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := rc.validate.Struct(response); err != nil {
		fmt.Printf("Validation failed! [%s]\n", err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Insert date
	utils.InsertResponseDate(&response, constants.ACTION_DATE_CREATE)
	// Create
	err := rc.ResponseService.CreateResponse(&response)
	if err != nil {
		fmt.Printf("Failed to create object! [%s]", err.Error())
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (rc *ResponseController) GetAll(ctx *gin.Context) {
	fmt.Printf("Getting response!\n")
	// Check admin
	admErr := utils.CheckAdmin(ctx, rc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	responses, err := rc.ResponseService.GetAll()
	if err != nil {
		fmt.Printf("Failed to get object! [%s]", err.Error())
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, responses)
}

func (rc *ResponseController) GetByToken(ctx *gin.Context) {
	fmt.Printf("Getting response!\n")
	// Check admin
	admErr := utils.CheckAdmin(ctx, rc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	// Check token
	token := ctx.Query("token")
	if token == "" {
		customErr.ThrowCustomError(&customErr.ErrQueryParamMissing, ctx)
		return
	}
	responses, err := rc.ResponseService.GetByToken(&token)
	if err != nil {
		fmt.Printf("Failed to get object! [%s]", err.Error())
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, responses)
}

func (rc *ResponseController) BatchDeleteResponse(ctx *gin.Context) {
	fmt.Printf("Deleting response!\n")
	admErr := utils.CheckAdmin(ctx, rc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	token := ctx.Query("token")
	if token == "" {
		customErr.ThrowCustomError(&customErr.ErrQueryParamMissing, ctx)
		return
	}
	err := rc.ResponseService.BatchDeleteByToken(&token)
	if err != nil {
		fmt.Printf("Failed to delete object! [%s]", err.Error())
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (rc *ResponseController) RegisterResponseRoutes(rg *gin.RouterGroup) {
	responseroute := rg.Group("/responses")
	responseroute.POST("/create", rc.CreateResponse)
	responseroute.GET("/getall", rc.GetAll)
	responseroute.GET("/getbytoken", rc.GetByToken)
	responseroute.DELETE("/batchdelete", rc.BatchDeleteResponse)
}
