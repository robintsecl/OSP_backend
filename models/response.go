package models

import "time"

type Response struct {
	Token          string           `json:"token" bson:"token"`
	CreateDate     time.Time        `json:"createDate" bson:"createDate"`
	UpdateDate     time.Time        `json:"updateDate" bson:"updateDate"`
	ResponseAnswer []ResponseAnswer `json:"responseAnswer" bson:"responseAnswer"`
}

type ResponseAnswer struct {
	Title  string `json: "title" bson:"title"`
	Answer string `json: "answer" bson:"answer"`
}
