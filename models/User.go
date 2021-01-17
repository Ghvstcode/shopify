package models

import (
	"context"
	"github.com/GhvstCode/shopify-challenge/database"
	"github.com/GhvstCode/shopify-challenge/utils"
	l "github.com/GhvstCode/shopify-challenge/utils/logger"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

//Token struct represents a JWT Claim.
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
	ImageLinks []Image         `bson:"imageLink, omitempty" json:"imageLink"`
}

//containsAny is a simple utility function for checking if the password contains any of the reserved keywords.
func containsAny (s string, substr []string) bool{
	for _, v := range substr {
		if strings.Contains(s, v){
			return false
		}
	}
	return true
}

//genAuthToken generates a JWT token based on the current time/User ID.
func genAuthToken(id string) (string, error) {
	t := &Token{
		UserId: id + "_" + time.StampMicro,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(86400 * time.Minute).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), t)
	tokenString, err := token.SignedString([]byte("os.Getenv")) //Change this to load the jwt secret from env file.
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//validate is used to validate the input.
func (U *User) validate() *utils.Data {
	//Validate Email
	if !strings.Contains(U.Email, "@") {
		return utils.Response(false, "Provide a valid email address", http.StatusBadRequest)
	}
	//validate password length
	if len(U.Password) < 6 {
		return utils.Response(false, "Password is required", http.StatusBadRequest)
	}

	reserved := []string{"abcdefg","password", "blahblah", U.Email}
	if value := containsAny(U.Password,reserved); value != true {
		return utils.Response(false, "Please provide a valid password", http.StatusBadRequest)
	}

	ErrorChan := make(chan error, 1)
	defer close(ErrorChan)
	go func() {
		ErrorChan <- database.UserDB.FindOne(context.TODO(), bson.D{{"email",U.Email}}).Decode(U)
	}()
	Error := <-ErrorChan
	if Error == nil {
		return utils.Response(false, "Invalid Email!", http.StatusBadRequest)
	}

	return utils.Response(true, "Validated", http.StatusAccepted)
}

//Create is called by the handler and is used to create a new user.
func (U *User) Create() *utils.Data {
	resp := U.validate()
	ok := resp.Result
	if !ok {
		return resp
	}

	hashedPassword, err := utils.Hash(U.Password)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to create user", http.StatusInternalServerError)
	}
	U.Password = string(hashedPassword)

	res, err := database.UserDB.InsertOne(context.TODO(), &User{
		ID:       primitive.NewObjectID(),
		Email:    U.Email,
		Password: U.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to create user", http.StatusInternalServerError)
	}

	var UID string
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		UID = oid.Hex()
	}

	t, e := genAuthToken(UID)
	if e != nil {
		l.ErrorLogger.Println(e)
		return utils.Response(false, "Failed to create account, connection error.", http.StatusBadGateway)
	}

	U.Password = ""
	v := &UserResponse{
		ID:       UID,
		Email:    U.Email,
		Password: "",
	}
	response := utils.Response(true, "created", http.StatusCreated)
	response.Token = t
	response.Data = [1]*UserResponse{v}
	return response
}

func (U *User) Login() *utils.Data {
	user := &User{}
	errorChan := make(chan error, 1)
	defer close(errorChan)
	go func() {
		errorChan <- database.UserDB.FindOne(context.TODO(), bson.M{"email": U.Email, "verified": true}).Decode(user)
	}()
	err := <-errorChan
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Unable to  Login!", http.StatusBadRequest)
	}

	err = bcrypt.CompareHashAndPassword([]byte(U.Password), []byte(user.Password))


	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Invalid login credentials. Please try again", http.StatusUnauthorized)
	}

	// CHECK THE DATABASE FOR ALL THE LINKS BELONGING TO THIS USER AD RETURN IT

	cursor, err := database.PhotoDB.Find(context.TODO(), bson.M{"ownerId": user.ID})
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occured", http.StatusInternalServerError)
	}
	var results []Image
	if err = cursor.All(context.TODO(), &results); err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occured", http.StatusInternalServerError)
	}

	var UID string
	UID = user.ID.Hex()

	t, e := genAuthToken(UID)
	if e != nil {
		l.ErrorLogger.Println(e)
		return utils.Response(false, "Failed to Log in, connection error.", http.StatusBadGateway)
	}

	user.Password = ""
	v := &UserResponse{
		ID:       UID,
		Email:    user.Email,
		Password: "",
		ImageLinks: results,
	}
	response := utils.Response(true, "Logged In in", http.StatusCreated)
	response.Token = t
	response.Data = [1]*UserResponse{v}

	return response
}