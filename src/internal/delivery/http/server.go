package http

import (
	"github.com/gin-gonic/gin"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/http/middlewares"
	"github.com/idzharbae/digital-wallet/src/internal/gateway/rmq_gateway"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type HttpServer struct {
	engine      *gin.Engine
	rmqProducer *rmq_gateway.RMQProducer
	dbConn      *pgxpool.Pool
	redisClient *redis.Client

	userUC usecase.UserUC
}

func NewServer(rmqProducer *rmq_gateway.RMQProducer, dbConn *pgxpool.Pool, redisClient *redis.Client, userUC usecase.UserUC) *HttpServer {
	log.Info().Msg("Initializing service")

	// Create barebone engine
	engine := gin.New()
	// Add default recovery middleware
	engine.Use(gin.Recovery())

	// disabling the trusted proxy feature
	engine.SetTrustedProxies(nil)

	// Add cors, request ID and request logging middleware
	log.Info().Msg("Adding cors, request id and request logging middleware")
	engine.Use(middlewares.CORSMiddleware(), middlewares.RequestID(), middlewares.RequestLogger())

	server := &HttpServer{
		engine:      engine,
		rmqProducer: rmqProducer,
		dbConn:      dbConn,
		redisClient: redisClient,
		userUC:      userUC,
	}

	// Setup routers
	log.Info().Msg("Setting up routers")
	server.SetupRouters()

	return server
}

func (s *HttpServer) Run(addr ...string) error {
	return s.engine.Run(addr...)
}
