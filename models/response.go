package models

type Response struct {
	Token          string           `json: "token" bson:"token"`
	ResponseAnswer []ResponseAnswer `json: "responseAnswer" bson:"responseAnswer"`
}

type ResponseAnswer struct {
	Title  string `json: "title" bson:"title"`
	Answer string `json: "answer" bson:"answer"`
}
