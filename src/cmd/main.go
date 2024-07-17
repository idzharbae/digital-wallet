package main

import (
	"fmt"
	"os"

	"github.com/idzharbae/digital-wallet/src/internal/app"
	"github.com/idzharbae/digital-wallet/src/internal/utils"
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"
)

func main() {
	errChan := make(chan error)

	rmqConnectionString := utils.GetEnvVar("RMQ_URL")
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		utils.GetEnvVar("DB_USER"),
		utils.GetEnvVar("DB_PASS"),
		utils.GetEnvVar("DB_HOST"),
		utils.GetEnvVar("DB_PORT"),
		utils.GetEnvVar("DB_NAME"),
	)

	app, err := app.SetupApp(app.AppConf{
		RmqConnectionString: rmqConnectionString,
		DbConnectionString:  dbUrl,
	})
	if err != nil {
		panic(err)
	}

	// Run RabbitMQ consumer
	go func() {
		err := app.RabbitMq.Consume()
		if err != nil {
			errChan <- stacktrace.Propagate(err, "error occured when setting up RabbitMQ consumer")
		}
	}()

	// Run httpServer
	go func() {
		// Read ADDR and port
		addr := utils.GetEnvVar("GIN_ADDR")
		port := utils.GetEnvVar("GIN_PORT")
		// HTTP mode
		log.Info().Msgf("Starting service on http//:%s:%s", addr, port)
		if err := app.Http.Run(fmt.Sprintf("%s:%s", addr, port)); err != nil {
			errChan <- stacktrace.Propagate(err, "error occurred while setting up the http server")
		}
	}()

	app.Cron.Start()

	for err := range errChan {
		log.Error().Err(err).Msg("received error from one of the delivery layer")
		os.Exit(1)
	}
}
