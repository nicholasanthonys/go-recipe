package main

import (
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/nicholasanthonys/go-recipe/infrastructure"
	"github.com/nicholasanthonys/go-recipe/infrastructure/database"
	"github.com/nicholasanthonys/go-recipe/infrastructure/log"
	"github.com/nicholasanthonys/go-recipe/infrastructure/router"
	"github.com/nicholasanthonys/go-recipe/infrastructure/validation"
)

func main() {
	godotenv.Load(".env")

	var app = infrastructure.NewConfig().
		Name(os.Getenv("APP_NAME")).
		ContextTimeout(10 * time.Second).
		Logger(log.InstanceZapLogger).
		Validator(validation.InstanceGoPlayground).
		DbNoSQL(database.InstanceMongoDB)
		// DbSQL(database.InstancePostgres).

	app.WebServerPort(os.Getenv("APP_PORT")).
		WebServer(router.InstanceGin).
		Start()
}
