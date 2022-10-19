package main

import (
	"os"
	"time"

	"github.com/gsabadini/go-bank-transfer/infrastructure"
	"github.com/gsabadini/go-bank-transfer/infrastructure/log"
	"github.com/gsabadini/go-bank-transfer/infrastructure/router"
	"github.com/gsabadini/go-bank-transfer/infrastructure/validation"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	var app = infrastructure.NewConfig().
		Name(os.Getenv("APP_NAME")).
		ContextTimeout(10 * time.Second).
		Logger(log.InstanceZapLogger).
		Validator(validation.InstanceGoPlayground)
		// DbSQL(database.InstancePostgres).
		// DbNoSQL(database.InstanceMongoDB)

	app.WebServerPort(os.Getenv("APP_PORT")).
		WebServer(router.InstanceGin).
		Start()
}
