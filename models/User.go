package user

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Token struct {
	UserId string
	jwt.StandardClaims
}

type User struct {
ID primitive.ObjectID `bson:"_id, omitempty" json:"_id, omitempty"`
Email string `bson:"email" json:"email"`
Password string `bson:"password" json:"password"`
CreatedAt time.Time `bson:”created_at” json:”created_at,omitempty” `
UpdatedAt time.Time `json:”updated_at,omitempty” bson:”updated_at”`
}

type UserResponse struct {
	ID string `bson:"_id, omitempty" json:"_id, omitempty"`
	Email string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
	//Links []UserLinkRequest         `bson:"link, omitempty" json:"Link"`
}