package controllers

import (
	"github.com/GhvstCode/shopify-challenge/models"
	"github.com/GhvstCode/shopify-challenge/utils"
	l "github.com/GhvstCode/shopify-challenge/utils/logger"
	"net/http"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		l.ErrorLogger.Println(err)
		utils.Response(false, "Invalid request", http.StatusBadRequest).Send(w)
		return
	}

	ownerID := r.Context().Value("user")
	if ownerID == nil {
		utils.Response(false, "UnAuthorized Access", http.StatusUnauthorized).Send(w)
		return
	}

	res := models.Upload(file, fileHeader, ownerID.(string))
	res.Send(w)
}
