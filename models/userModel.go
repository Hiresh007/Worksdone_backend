package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      *string            `json:"username" validate:"required,min=2,max=100"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Token         *string            `json:"token"`
	Refresh_token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}

type ResumeScan struct {
	Role   *string `json:"role" validate:"required"`
	Url    *string `json:"url" validate:"required"`
	Prompt *string `json:"prompt"`
}
type VideoIdRequest struct {
	VideoUrl string `json:"videoUrl"`
}
type ATSScore struct {
	ResumePath  string `json:"resumeUrl,omitempty" validate:"required"`
	Description string `json:"description" validate:"required"`
}
