package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tpk32/academic-tracker-server/db"
	"github.com/tpk32/academic-tracker-server/router"
	"github.com/tpk32/academic-tracker-server/services"
)

type Config struct{
	Port string
}

type Application struct{
	Config Config
	Models services.Models
}

func(app *Application) Serve() error{
	//err := godotenv.Load() //dont need in railway
	// if(err != nil){
	// 	log.Fatal("Error laoding .env file while serving")
	// }
	port := os.Getenv("PORT")
	fmt.Println("API is running on port", port)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: router.Routes(),
	}
	return srv.ListenAndServe()
}

func main(){
	err := godotenv.Load()
	if(err != nil){
		log.Fatal("Error loading .env file in main")
	}

	cfg := Config{
		Port: os.Getenv("PORT"),
	}

	dsn := os.Getenv("DSN")
	dbConn, err := db.ConnectPostgres(dsn)
	if err != nil{
		log.Fatal("Cannot connect to database")
	}

	defer dbConn.DB.Close()

	app := &Application{
		Config: cfg,
		Models: services.New(dbConn.DB),
	}

	err = app.Serve()
	if err != nil {
		log.Fatal(err)
	}
}