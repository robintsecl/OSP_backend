package utils

import (
	"github.com/gin-gonic/gin"
	constant "github.com/robintsecl/osp_backend/constants"
	customErr "github.com/robintsecl/osp_backend/errors"
	"github.com/robintsecl/osp_backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func CheckAdmin(ctx *gin.Context, usercollection *mongo.Collection) error {
	name := ctx.Query("name")
	pw := ctx.Query("password")
	if name == "" || pw == "" {
		return customErr.ErrUnauthorized
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
	// Check duplicate title
	titleMap := make(map[string]bool)
	isDupTitle := false
	// Check format and spec
	isWrongFormat := false
	for _, question := range *questions {
		checkIsDuplicateTitle(question.Title, titleMap, &isDupTitle)
		checkFormatAndSpec(question, &isWrongFormat)
	}
	if isDupTitle {
		return customErr.ErrDuplicateQuestionTitle
	}
	if isWrongFormat {
		return customErr.ErrInvalidQuestionFormatAndSpec
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

func checkFormatAndSpec(question models.Question, isWrongFormat *bool) {
	if !*isWrongFormat {
		if question.Type == constant.TEXTBOX {
			if len(question.Spec) > constant.TEXTBOX_MIN_LEN {
				*isWrongFormat = true
			}
		}
		if question.Type == constant.MC {
			if len(question.Spec) < constant.MC_MIN_LEN {
				*isWrongFormat = true
			}
		}
		if question.Type == constant.LS {
			if len(question.Spec) < constant.LS_MIN_LEN {
				*isWrongFormat = true
			}
		}
	}
}
