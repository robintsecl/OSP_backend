package models

import "time"

type Survey struct {
	Title      string     `json:"title" validate:"required,min=1,max=5" bson:"title"`
	Token      string     `json:"token" bson:"token"`
	CreateDate time.Time  `json:"createDate" bson:"createDate,omitempty"`
	UpdateDate time.Time  `json:"updateDate" bson:"updateDate,omitempty"`
	Questions  []Question `json:"questions" validate:"required,min=1,max=255" bson:"questions"`
}

type Question struct {
	Title string   `json:"title" validate:"required,min=1,max=2" bson:"title"`
	Type  string   `json:"type" validate:"required" bson:"type"`
	Spec  []string `json:"spec" validate:"dive,min=1,max=255" bson:"spec"`
}
