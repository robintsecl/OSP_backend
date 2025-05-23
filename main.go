package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/robintsecl/osp_backend/controllers"
	"github.com/robintsecl/osp_backend/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server             *gin.Engine
	surveyservice      services.SurveyService
	surveycontroller   controllers.SurveyController
	surveycollection   *mongo.Collection
	responseservice    services.ResponseService
	responsecontroller controllers.ResponseController
	responsecollection *mongo.Collection
	ctx                context.Context
	mongoclient        *mongo.Client
	err                error
)

func init() {
	ctx = context.TODO()
	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("mongodb is connected")
	// Init a admin user for some admin action
	usercollection := mongoclient.Database("osp-db").Collection("user")
	err := usercollection.FindOne(ctx, bson.D{bson.E{Key: "name", Value: "admin"}}).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			usercollection.InsertOne(ctx, bson.D{
				{Key: "name", Value: "admin"},
				{Key: "password", Value: "hkuabc123"},
			})
		} else {
			log.Fatal(err)
		}
	}
	fmt.Println("admin password: hkuabc123")

	validate := validator.New()

	surveycollection = mongoclient.Database("osp-db").Collection("survey")
	surveyservice = services.NewSurveyService(surveycollection, ctx)
	surveycontroller = controllers.NewSurveyController(surveyservice, usercollection, validate)

	responsecollection = mongoclient.Database("osp-db").Collection("response")
	responseservice = services.NewResponseService(responsecollection, surveycollection, ctx)
	responsecontroller = controllers.NewResponseController(responseservice, usercollection, validate)

	gin.SetMode(gin.ReleaseMode)
	server = gin.New()
}

func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/osp")
	surveycontroller.RegisterSurveyRoutes(basepath)
	responsecontroller.RegisterResponseRoutes(basepath)
	log.Fatal(server.Run(":9091"))
}
