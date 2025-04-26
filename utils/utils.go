package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robintsecl/osp_backend/constants"
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
		fmt.Printf(("Name or password is empty in query parameter!\n"))
		return customErr.ErrUnauthorized
	}
	query := bson.D{bson.E{Key: "name", Value: name}, bson.E{Key: "password", Value: pw}}
	err := usercollection.FindOne(ctx, query).Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Printf(("Name or password not matched!\n"))
			return customErr.ErrUnauthorized
		} else {
			return err
		}
	}
	fmt.Printf(("Login succeed!\n"))
	return nil
}

func CommonChecking(questions []*models.Question) error {
	// Check duplicate title
	titleMap := make(map[string]bool)
	isDupTitle := false
	// Check format and spec
	isWrongFormat := false
	for _, question := range questions {
		checkIsDuplicateTitle(question.Title, titleMap, &isDupTitle)
		checkFormatAndSpec(question, &isWrongFormat)
	}
	if isDupTitle {
		fmt.Printf("Duplicate title found!\n")
		return customErr.ErrDuplicateQuestionTitle
	}
	if isWrongFormat {
		fmt.Printf("Question format and specification is not matched!\n")
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
		fmt.Printf("Title [%s] is duplicated!\n", currentValue)
		*isDup = true
	}
}

func checkFormatAndSpec(question *models.Question, isWrongFormat *bool) {
	if *isWrongFormat {
		return
	}
	switch question.Type {
	case constant.TEXTBOX:
		if len(question.Spec) < constant.TEXTBOX_MIN_LEN {
			fmt.Printf("Length of [%s] specification is smaller than [%v]!\n", constant.TEXTBOX, constant.TEXTBOX_MIN_LEN)
			*isWrongFormat = true
		}
	case constant.MC:
		if len(question.Spec) < constant.MC_MIN_LEN {
			fmt.Printf("Length of [%s] specification is smaller than [%v]!\n", constant.MC, constant.MC_MIN_LEN)
			*isWrongFormat = true
		}
	case constant.LS:
		if len(question.Spec) < constant.LS_MIN_LEN {
			fmt.Printf("Length of [%s] specification is smaller than [%v]!\n", constant.LS, constant.LS_MIN_LEN)
			*isWrongFormat = true
		}
	default:
		fmt.Printf("Unknown question format detected!\n")
		*isWrongFormat = true
	}
}

func ResponseInputChecking(questions []*models.Question, answers []*models.ResponseAnswer) error {
	questionMap := make(map[string]models.Question)
	// Put question to map object
	for _, question := range questions {
		questionMap[question.Title] = *question
	}

	// Loop through answer
	for _, answer := range answers {
		question, isTitleExists := questionMap[answer.Title]
		// If answer title doesn't exists in question of that survey, throw error
		if !isTitleExists {
			fmt.Printf("Title [%s] in answer does not exists in question!\n", answer.Title)
			return customErr.ErrTitleNotFoundInQuestion
		}
		if question.Type == constant.TEXTBOX {
			// If user haven't input anything, throw error
			if len(answer.Answer) < 1 {
				fmt.Printf("Textbox answer is empty!\n")
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
				fmt.Printf("Answer [%s] is not selectable in question specification!\n", answer.Answer)
				return customErr.ErrInvalidAnswerInSpec
			}
		}
	}
	return nil
}

func InsertSurveyDate(survey *models.Survey, insertType string) {
	if insertType == constants.ACTION_DATE_CREATE {
		survey.CreateDate = time.Now()
		survey.UpdateDate = time.Now()
	}
	if insertType == constants.ACTION_DATE_UPDATE {
		survey.UpdateDate = time.Now()
	}
}

func InsertResponseDate(response *models.Response, insertType string) {
	if insertType == constants.ACTION_DATE_CREATE {
		response.CreateDate = time.Now()
		response.UpdateDate = time.Now()
	}
	if insertType == constants.ACTION_DATE_UPDATE {
		response.UpdateDate = time.Now()
	}
}
