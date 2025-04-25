package models

type Survey struct {
	Title     string     `json: "title" bson:"title"`
	Token     string     `json: "token" bson:"token"`
	Questions []Question `json: "questions" bson:"questions"`
}

type Question struct {
	Title string   `json: "title" bson:"title"`
	Type  string   `json: "type" bson:"type"`
	Spec  []string `json: "spec" bson:"spec"`
}
