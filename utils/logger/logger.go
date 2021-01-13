package logger

import (
	"io"
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)

	InfoLogger = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(mw, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}