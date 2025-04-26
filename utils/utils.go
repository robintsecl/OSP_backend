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
	if *isDup {
		return
	}
	if _, dup := titleMap[currentValue]; !dup {
		titleMap[currentValue] = true
	} else {
		*isDup = true
	}
}

func checkFormatAndSpec(question models.Question, isWrongFormat *bool) {
	if *isWrongFormat {
		return
	}
	switch question.Type {
	case constant.TEXTBOX:
		if len(question.Spec) < constant.TEXTBOX_MIN_LEN {
			*isWrongFormat = true
		}
	case constant.MC:
		if len(question.Spec) < constant.MC_MIN_LEN {
			*isWrongFormat = true
		}
	case constant.LS:
		if len(question.Spec) < constant.LS_MIN_LEN {
			*isWrongFormat = true
		}
	default:
		*isWrongFormat = true
	}
}

func ResponseInputChecking(questions *[]models.Question, answers *[]models.ResponseAnswer) error {
	questionMap := make(map[string]models.Question)
	// Put question to map object
	for _, question := range *questions {
		questionMap[question.Title] = question
	}

	// Loop through answer
	for _, answer := range *answers {
		question, isTitleExists := questionMap[answer.Title]
		// If answer title doesn't exists in question of that survey, throw error
		if !isTitleExists {
			return customErr.ErrTitleNotFoundInQuestion
		}
		if question.Type == constant.TEXTBOX {
			// If user haven't input anything, throw error
			if len(answer.Answer) < 1 {
				return customErr.ErrTextAnswerIsEmpty
			}
		}
		if question.Type == constant.LS || question.Type == constant.MC {
			isSpec := false
			// Loop through spec to see if user selected answer in spec
			for _, spec := range question.Spec {
				if isSpec {
					break
				}
				if spec == answer.Answer {
					isSpec = true
				}
			}
			if !isSpec {
				return customErr.ErrInvalidAnswerInSpec
			}
		}
	}
	return nil
}
