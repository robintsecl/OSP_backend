package models

import "time"

type Response struct {
	Token          string            `json:"token" validate:"required,min=1,max=5" bson:"token"`
	CreateDate     time.Time         `json:"createDate" bson:"createDate"`
	UpdateDate     time.Time         `json:"updateDate" bson:"updateDate"`
	ResponseAnswer []*ResponseAnswer `json:"responseAnswer" validate:"required,min=1,max=255,dive,required" bson:"responseAnswer"`
}

type ResponseAnswer struct {
	Title  string `json: "title" validate:"required,min=1,max=255" bson:"title"`
	Answer string `json: "answer" validate:"required,min=1,max=255" bson:"answer"`
}
