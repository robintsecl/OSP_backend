package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server           *gin.Engine
	ctx              context.Context
	surveycollection *mongo.Collection
	mongoclient      *mongo.Client
	err              error
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
	surveycollection = mongoclient.Database("osp-db").Collection("survey")
	server = gin.Default()
}

func main() {
	defer mongoclient.Disconnect(ctx)

	log.Fatal(server.Run(":9091"))
}
