package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckAdmin(ctx *gin.Context, usercollection *mongo.Collection) error {
	name := ctx.Query("name")
	pw := ctx.Query("password")
	if name == "" || pw == "" {
		return fmt.Errorf("name or pw query parameter is required")
	}
	query := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "password", Value: pw}}
	err := usercollection.FindOne(ctx, query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("incorrect name or password")
		} else {
			return err
		}
	}
	return nil
}
