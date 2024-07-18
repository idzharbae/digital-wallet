package http

import (
	"github.com/gin-gonic/gin"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/http/middlewares"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	"github.com/rs/zerolog/log"
)

type HttpServer struct {
	engine *gin.Engine

	userUC        usecase.UserUC
	transactionUC usecase.TransactionUC
}

func NewServer(userUC usecase.UserUC, transactionUC usecase.TransactionUC) *HttpServer {
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
		engine:        engine,
		userUC:        userUC,
		transactionUC: transactionUC,
	}

	// Setup routers
	log.Info().Msg("Setting up routers")
	server.SetupRouters()

	return server
}

func (s *HttpServer) Run(addr ...string) error {
	return s.engine.Run(addr...)
}
