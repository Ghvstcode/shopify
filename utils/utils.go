package utils

import (
	"encoding/json"
	"net/http"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	l "github.com/GhvstCode/shopify-challenge/utils/logger"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		l.ErrorLogger.Println("No .env file found")
	}
}

type Data struct {
	StatusCode int
	Message    string
	Result     bool
	Data       interface{}
	Token      string `json:"Token,omitempty"`
}

func (data Data) Send(w http.ResponseWriter) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(data.StatusCode)
	if err := json.NewEncoder(w).Encode(data);  err != nil  {
		return err
	}
	return nil
}

//Message is exported
func Message(result bool, message string) map[string]interface{} {
	return map[string]interface{}{"result": result, "message": message}
}

//Hash function encrypts the password. It uses the bcrypt library under the hood.
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func Response(result bool, message string, statusCode int) *Data {
	return &Data{
		StatusCode: statusCode,
		Message:    message,
		Result:     result,
	}
}





