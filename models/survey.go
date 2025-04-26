package models

import "time"

type Survey struct {
	Title      string      `json:"title" validate:"required,min=1,max=255" bson:"title"`
	Token      string      `json:"token" bson:"token"`
	CreateDate time.Time   `json:"createDate" bson:"createDate"`
	UpdateDate time.Time   `json:"updateDate" bson:"updateDate"`
	Questions  []*Question `json:"questions" validate:"required,min=1,max=255,dive,required" bson:"questions"`
}

type Question struct {
	Title string   `json:"title" validate:"required,min=1,max=255" bson:"title"`
	Type  string   `json:"type" validate:"required" bson:"type"`
	Spec  []string `json:"spec" validate:"required,dive,min=1,max=255" bson:"spec"`
}
