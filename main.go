package main

import (
	"Engine/api"
	"Engine/storage"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	Port       string
	Enviroment string
	DBConn      string
)

func main() {
	// set up environment variables
	SetupEnviroment()

	// connect to database
	db , err := storage.NewDB(DBConn)
	if err!= nil {
		logrus.Fatal(err)
	}

    logrus.Info("Established a successful database connection.")

	// Initialize handlers
	r := chi.NewRouter()
	api.InitHandlers(r,db)

}

func SetupEnviroment() {
	//!!! need to be commented out if pushing to dev
	 err := godotenv.Load()
	 if err != nil {
	 	panic("Error loading .env file")
	 }
	//!!! need to be commented out if pushing to dev

	Enviroment = os.Getenv("ENV")
	if Enviroment == "" {
		panic("Enviroment enviroment vairable must be set")
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		panic("Port enviroment vairable must be set")
	}

	DBConn = os.Getenv("DBCONN")
	if DBConn == "" {
		panic("DBCONN enviroment vairable must be set")
	}

	logrus.Info("enviroment set...")
}
