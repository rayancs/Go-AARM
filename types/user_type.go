package types

type UserType struct {
	Id         string `json:"id" bson:"_id,omitempty"`
	Email      string `json:"email" bson:"email,omitempty"`
	Name       string `json:"name" bson:"name"`
	Emoji      string `json:"emoji" bson:"emoji"`
	PictureURL string `json:"pictureURL" bson:"pictureURL"`
}
type UserResponseType struct {
	Emoji      string `json:"emoji" bson:"emoji"`
	Name       string `json:"name" bson:"name"`
	Email      string `json:"email" bson:"email,omitempty"`
	PictureURL string `json:"pictureURL" bson:"pictureURL"`
}
type UserDbResponseType struct {
	Email      string `json:"email" bson:"email,omitempty"`
	Id         string `json:"id" bson:"_id,omitempty"`
	Name       string `json:"name" bson:"name"`
	Emoji      string `json:"emoji" bson:"emoji"`
	PictureURL string `json:"pictureURL" bson:"pictureURL"`
}

// to find
type UserRequestType struct {
	Email      string `json:"email" bson:"email,omitempty"`
	Emoji      string `json:"emoji" bson:"emoji"`
	Name       string `json:"name" bson:"name"`
	PictureURL string `json:"pictureURL" bson:"pictureURL"`
}
