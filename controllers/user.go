package controllers

import (
	"encoding/json"
	"github.com/GhvstCode/shopify-challenge/models"
	"github.com/GhvstCode/shopify-challenge/utils"
	l "github.com/GhvstCode/shopify-challenge/utils/logger"
	"net/http"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		l.ErrorLogger.Println(err)
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}
	res := user.Create()
	res.Send(w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		l.ErrorLogger.Println(err)
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}


	res := user.Login()
	res.Send(w)
}

