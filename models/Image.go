package models

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/GhvstCode/shopify-challenge/database"
	"github.com/GhvstCode/shopify-challenge/utils"
	l "github.com/GhvstCode/shopify-challenge/utils/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Image struct {
	ID primitive.ObjectID `bson:"_id, omitempty" json:"_id,omitempty"`
	FileName string `bson:"fileName, omitempty" json:"fileName"`
	Url string `bson:"url, omitempty" json:"url"`
	OwnerID primitive.ObjectID `bson:"ownerId, omitempty" json:"ownerId"`
	Size int64 `bson:"size, omitempty" json:"size"`
	CreatedAt time.Time `json:”created_at,omitempty” bson:”created_at”`
}

type ImageResponse struct {
	ID string `bson:"_id, omitempty" json:"_id,omitempty"`
	FileName string `bson:"fileName, omitempty" json:"fileName"`
	Url string `bson:"url, omitempty" json:"url"`
	Size int64 `bson:"size, omitempty" json:"size"`
	CreatedAt time.Time `json:”created_at,omitempty” bson:”created_at”`
}

func Upload(file multipart.File, fileHeader *multipart.FileHeader, id string)*utils.Data{
 	//Validate File to be uploaded to see if they have the right extensions.
	bucket := os.Getenv("CLOUD_STORAGE_BUCKET_NAME")
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Error creating Storage Client", http.StatusInternalServerError)
	}

	defer file.Close()

	fileObject := storageClient.Bucket(bucket).Object(fileHeader.Filename)

	sw := fileObject.NewWriter(ctx)



	if _, err := io.Copy(sw, file); err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Error creating Storage Client", http.StatusInternalServerError)

	}

	if err := sw.Close(); err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Error creating Storage Client", http.StatusInternalServerError)

	}

	fileAcl := fileObject.ACL()
	if err := fileAcl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Error creating Storage Client", http.StatusInternalServerError)
	}
	//fmt.Println(sw.Attrs().MediaLink)
	fileAttr, err := fileObject.Attrs(ctx)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "Error creating Storage Client", http.StatusInternalServerError)
	}


	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An Error occurred, Unauthourized Access", http.StatusInternalServerError)
	}

	res, err := database.PhotoDB.InsertOne(context.TODO(), &Image{
		ID:       primitive.NewObjectID(),
		FileName:  fileHeader.Filename,
		Url:       fileAttr.MediaLink,
		OwnerID:   oid,
		Size:      fileAttr.Size,
		CreatedAt: time.Now(),
	})

	if err != nil {
		l.ErrorLogger.Println(err)
		return utils.Response(false, "An error occurred! Unable to Shorten Link", http.StatusInternalServerError)
	}

	var UID string
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		UID = oid.Hex()
	}



	v := &ImageResponse{
		ID:        UID,
		FileName:  fileHeader.Filename,
		Url:       fileAttr.MediaLink,
		Size:      fileAttr.Size,
		CreatedAt: time.Now(),
	}

	response := utils.Response(true, "created", http.StatusCreated)
	response.Data = [1]*ImageResponse{v}
	return response

}