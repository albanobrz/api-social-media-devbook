package entities

type Password struct {
	New     string `json:"new" bson:"new"`
	Current string `json:"current" bson:"current"`
}
