package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrUnauthorized                 = errors.New("unauthorized access")
	ErrDuplicateQuestionTitle       = errors.New("duplicate question title")
	ErrInvalidQuestionFormatAndSpec = errors.New("invalid question specification according to format")
	ErrQueryParamMissing            = errors.New("query parameter missing")
	ErrDataNotFound                 = errors.New("data not found")
	ErrDBBadGateway                 = errors.New("bad gateway")
	ErrInvalidQuestionFormat        = errors.New("invalid question format and specification")
	ErrTitleNotFoundInQuestion      = errors.New("title of answer not found in question")
	ErrTextAnswerIsEmpty            = errors.New("answer is empty")
	ErrInvalidAnswerInSpec          = errors.New("invalid answer according to specification")
	// ErrSurveyInsertionFailed  = errors.New("failed to insert survey")
	// ErrInvalidSurveyData      = errors.New("invalid survey data")
)

func ThrowCustomError(err *error, ctx *gin.Context) {
	switch *err {
	case ErrUnauthorized:
		ctx.JSON(http.StatusForbidden, gin.H{"message": ErrUnauthorized.Error()})
	case ErrDuplicateQuestionTitle:
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrDuplicateQuestionTitle.Error()})
	case ErrInvalidQuestionFormatAndSpec:
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrInvalidQuestionFormatAndSpec.Error()})
	case ErrQueryParamMissing:
		ctx.JSON(http.StatusBadRequest, gin.H{"message": ErrQueryParamMissing.Error()})
	case ErrDataNotFound:
		ctx.JSON(http.StatusNotFound, gin.H{"message": ErrDataNotFound.Error()})
	case ErrDBBadGateway:
		ctx.JSON(http.StatusBadGateway, gin.H{"message": ErrDBBadGateway.Error()})
	case ErrTitleNotFoundInQuestion:
		ctx.JSON(http.StatusBadGateway, gin.H{"message": ErrTitleNotFoundInQuestion.Error()})
	case ErrTextAnswerIsEmpty:
		ctx.JSON(http.StatusBadGateway, gin.H{"message": ErrTextAnswerIsEmpty.Error()})
	case ErrInvalidAnswerInSpec:
		ctx.JSON(http.StatusBadGateway, gin.H{"message": ErrInvalidAnswerInSpec.Error()})
	default:
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err})
	}
}
