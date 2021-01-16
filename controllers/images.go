package controllers

import (
	"fmt"
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

	fmt.Println(fileHeader.Filename)
	res := models.Upload(file, fileHeader)
	res.Send(w)
}
