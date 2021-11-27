package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"trendyolcase/pkg/api"
	"trendyolcase/pkg/repository/link"
	"trendyolcase/pkg/service"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	a := App{}
	_ = a.initialize(os.Getenv("CONNECTION_STRING"))
	a.routes()
	err = a.run(os.Getenv("SERV_PORT"))
}

func (a *App) initialize(connectionAdrr string) error {
	var err error
	connectionString := fmt.Sprintf("%s", connectionAdrr)
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
		return err
	}
	a.Router = mux.NewRouter()
	return err
}

func (a *App) run(addr string) error {
	err := http.ListenAndServe(addr, a.Router)
	if err != nil {
		log.Fatalf("Port is using already.")
		return err
	}
	return err
}

func (a *App) routes() {
	converterAPI := InitConverterAPI(a.DB)
	a.Router.HandleFunc("/getDeepLink", converterAPI.GenerateDeepLink()).Methods("POST")
	a.Router.HandleFunc("/getWebURL", converterAPI.GenerateWebURL()).Methods("POST")
}

func InitConverterAPI(db *sql.DB) api.ConverterAPI {
	converterRepository := link.NewRepository(db)
	converterService := service.NewConverterService(converterRepository)
	converterAPI := api.NewConverterAPI(converterService)
	return converterAPI
}
