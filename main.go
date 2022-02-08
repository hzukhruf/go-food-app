package main

import (
	"fmt"
	"food-app/api"
	"food-app/pkg/config"
	"food-app/pkg/database"
	"log"

	"github.com/spf13/viper"
)


func init() {
	config.GetConfig()
}

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	port := fmt.Sprintf(":%d",viper.GetInt("App.Port"))

	app := api.SetupRouter(db)
	app.Run(port)
}