package main

import (
	"log"

	"github.com/nirrax/url_shortener/cmd/api"
	"github.com/nirrax/url_shortener/internal/config"
	"github.com/nirrax/url_shortener/internal/database/migrations"
	"github.com/nirrax/url_shortener/internal/database/postgres"
	"github.com/nirrax/url_shortener/internal/repository"
	"github.com/nirrax/url_shortener/internal/service"
	"github.com/nirrax/url_shortener/internal/utils"
)

func main() {
	conf := config.LoadConfig()
	log.Println("Config loaded")

	db, err := postgres.NewPostgresDB(conf.PostgresConfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
	defer db.Close()

	mig := migrations.NewGooseMigration(db.DB.DB, conf)
	err = mig.Up()
	if err != nil {
		log.Println(err)
	}

	repository := repository.NewRepository(db)
	encoder := utils.NewBase62Encoder()
	service := service.NewService(repository, encoder)

	app := api.App{
		Config:  conf,
		Service: service,
	}

	log.Println("Server starts on port :", conf.ServerPort)
	log.Fatal(app.Run())
}
