package utils

import (
	"fmt"

	"github.com/gin-gonic/gin"
	customErr "github.com/robintsecl/osp_backend/errors"
	"github.com/robintsecl/osp_backend/models"
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
			return customErr.ErrUnauthorized
		} else {
			return err
		}
	}
	return nil
}

func CommonChecking(questions *[]models.Question) error {
	titleMap := make(map[string]bool)
	isDupTitle := false
	for _, question := range *questions {
		checkIsDuplicateTitle(question.Title, titleMap, &isDupTitle)
	}
	if isDupTitle {
		return customErr.ErrDuplicateQuestionTitle
	}
	return nil
}

func checkIsDuplicateTitle(currentValue string, titleMap map[string]bool, isDup *bool) {
	if !*isDup {
		if _, dup := titleMap[currentValue]; !dup {
			titleMap[currentValue] = true
		} else {
			*isDup = true
		}
	}
}
