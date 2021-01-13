package database

import (
"context"
"os"
_ "time"

"github.com/joho/godotenv"
"go.mongodb.org/mongo-driver/mongo"
"go.mongodb.org/mongo-driver/mongo/options"

l "github.com/GhvstCode/Shortly/utils/logger"
)
var (
	UserDB *mongo.Collection
	PhotoDB *mongo.Collection
)
var ctx = context.TODO()

func init() {
	if err := godotenv.Load(); err != nil {
		l.ErrorLogger.Println("No .env file found")
	}
	l.InfoLogger.Println("Connecting to DB...")
	envUri, ok := os.LookupEnv("MongoDB_URI")

	Uri := envUri

	if !ok {

		l.WarningLogger.Println("Unable to connect to load connection URI from env file,connecting to local db!")
		Uri = "mongodb://localhost:27017"
		//Uri = "mongodb://mongo:27017"
	}
	clientOptions := options.Client().ApplyURI(Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		l.ErrorLogger.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		l.ErrorLogger.Fatal(err)
	}

	LinkDB = client.Database("Shortly").Collection("link")
	UserDB = client.Database("Shortly").Collection("user")
}
