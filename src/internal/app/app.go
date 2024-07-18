package app

import (
	"fmt"

	"github.com/idzharbae/digital-wallet/src/internal/delivery/cronjob"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/http"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/rabbitmq"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/postgresql_gateway"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/redis_gateway"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/rmq_gateway"
	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	"github.com/idzharbae/digital-wallet/src/internal/utils"
	"github.com/palantir/stacktrace"
	"github.com/robfig/cron/v3"
)

type AppConf struct {
	RmqConnectionString string
	DbConnectionString  string
}

type DigitalWalletApp struct {
	RabbitMq *rabbitmq.RabbitMQConsumer
	Http     *http.HttpServer
	Cron     *cronjob.Cron
}

// Function to setup the app object
func SetupApp(conf AppConf) (*DigitalWalletApp, error) {
	rmqProducer, err := rmq_gateway.NewProducer(conf.RmqConnectionString)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error when setting up rabbit MQ producer")
	}

	rmqConsumerConn, err := rmq_gateway.NewConsumer(conf.RmqConnectionString)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error when setting up rabbit MQ consumer")
	}

	pgPool, err := postgresql_gateway.NewPgPool(conf.DbConnectionString)
	if err != nil {
		return nil, stacktrace.Propagate(err, "error when setting up pg pool")
	}

	redisClient := redis_gateway.NewRedisClient(
		fmt.Sprintf("%s:%s", utils.GetEnvVar("REDIS_ADDR"), utils.GetEnvVar("REDIS_PORT")),
		utils.GetEnvVar("REDIS_PASS"),
	)

	userTokenRepo := repository.NewUserToken(pgPool, redisClient)
	userBalanceRepo := repository.NewUserBalance(pgPool, redisClient)
	userTransactionRepo := repository.NewUserTransaction(pgPool, redisClient)
	transactionHandler := repository.NewTransactionHandler(pgPool)
	userUC := usecase.NewUser(userTokenRepo, userBalanceRepo, transactionHandler)
	transactionUC := usecase.NewTransaction(transactionHandler, userTransactionRepo, userBalanceRepo, rmqProducer)
	httpServer := http.NewServer(userUC, transactionUC)
	rmqConsumer := rabbitmq.NewConsumer(rmqConsumerConn, userTransactionRepo)
	cronJob := cronjob.NewCron(cron.New(), userTransactionRepo)

	return &DigitalWalletApp{
		RabbitMq: rmqConsumer,
		Http:     httpServer,
		Cron:     cronJob,
	}, nil
}
