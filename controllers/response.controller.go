package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	customErr "github.com/robintsecl/osp_backend/errors"
	"github.com/robintsecl/osp_backend/models"
	"github.com/robintsecl/osp_backend/services"
	utils "github.com/robintsecl/osp_backend/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type ResponseController struct {
	ResponseService services.ResponseService
	usercollection  *mongo.Collection
}

func NewResponseController(responseservice services.ResponseService, usercollection *mongo.Collection) ResponseController {
	return ResponseController{
		ResponseService: responseservice,
		usercollection:  usercollection,
	}
}

func (rc *ResponseController) CreateResponse(ctx *gin.Context) {
	// Bind json
	var response models.Response
	if err := ctx.ShouldBindJSON((&response)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	// Create
	err := rc.ResponseService.CreateResponse(&response)
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (rc *ResponseController) GetAll(ctx *gin.Context) {
	// Check admin
	admErr := utils.CheckAdmin(ctx, rc.usercollection)
	if admErr != nil {
		customErr.ThrowCustomError(&admErr, ctx)
		return
	}
	responses, err := rc.ResponseService.GetAll()
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, responses)
}

// TODO: pass token
func (rc *ResponseController) GetByToken(ctx *gin.Context) {
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
	responses, err := rc.ResponseService.GetAll()
	if err != nil {
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, responses)
}

func (rc *ResponseController) BatchDeleteResponse(ctx *gin.Context) {
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
		customErr.ThrowCustomError(&err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (rc *ResponseController) RegisterResponseRoutes(rg *gin.RouterGroup) {
	responseroute := rg.Group("/response")
	responseroute.POST("/create", rc.CreateResponse)
	responseroute.GET("/getall", rc.GetAll)
	responseroute.PUT("/getbytoken", rc.GetByToken)
	responseroute.DELETE("/batchdelete", rc.BatchDeleteResponse)
}
