package http

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/http/dto"
	"github.com/idzharbae/digital-wallet/src/internal/usecase"
	"github.com/idzharbae/digital-wallet/src/internal/utils"
	"github.com/palantir/stacktrace"
	"github.com/rs/zerolog/log"
)

func (s *HttpServer) RegisterUser(c *gin.Context) {
	requestId := c.GetString("x-request-id")
	var request dto.RegisterUserRequest
	if binderr := c.ShouldBindJSON(&request); binderr != nil {
		log.Error().Err(binderr).Str("requestId", requestId).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request json",
		})
		return
	}

	if len(request.Username) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username is empty",
		})
		return
	}

	if len(request.Username) > 256 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username is too long",
		})
		return
	}

	if !utils.ValidateUserName(request.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username must start with alphabet and must be alphanumeric or underscore",
		})
		return
	}

	token, err := s.userUC.RegisterUser(c.Request.Context(), request.Username)
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Msg("Failed to register user")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to register user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *HttpServer) BalanceRead(c *gin.Context) {
	requestId := c.GetString("x-request-id")

	userToken := c.GetHeader("Authorization")
	if len(userToken) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "header 'Authorization' is empty",
		})
		return
	}

	username, err := s.userUC.GetUserNameFromToken(c.Request.Context(), userToken)
	if stacktrace.RootCause(err) == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to get user by token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user information",
		})
		return
	}

	balance, err := s.userUC.GetUserBalance(c.Request.Context(), username)
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to get user balance")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user balance",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func (s *HttpServer) BalanceTopUp(c *gin.Context) {
	requestId := c.GetString("x-request-id")

	userToken := c.GetHeader("Authorization")
	if len(userToken) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "header 'Authorization' is empty",
		})
		return
	}

	username, err := s.userUC.GetUserNameFromToken(c.Request.Context(), userToken)
	if stacktrace.RootCause(err) == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to get user by token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user information",
		})
		return
	}

	var request dto.BalanceTopUpRequest
	if binderr := c.ShouldBindJSON(&request); binderr != nil {
		log.Error().Err(binderr).Str("requestId", requestId).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request json",
		})
		return
	}

	balance, err := s.userUC.TopUpUserBalance(c.Request.Context(), username, request.Amount)
	if err == usecase.ErrTopUpTooLarge {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to topup user balance")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to topup balance",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

func (s *HttpServer) Transfer(c *gin.Context) {
	requestId := c.GetString("x-request-id")

	userToken := c.GetHeader("Authorization")
	if len(userToken) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "header 'Authorization' is empty",
		})
		return
	}

	username, err := s.userUC.GetUserNameFromToken(c.Request.Context(), userToken)
	if stacktrace.RootCause(err) == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to get user by token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user information",
		})
		return
	}

	var request dto.TransferRequest
	if binderr := c.ShouldBindJSON(&request); binderr != nil {
		log.Error().Err(binderr).Str("requestId", requestId).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request json",
		})
		return
	}

	err = s.transactionUC.TransferBalance(c.Request.Context(), username, request.ToUsername, request.Amount)
	if err == usecase.ErrNotEnoughBalance {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to transfer user balance")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to transfer balance",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (s *HttpServer) TopUsers(c *gin.Context) {
	topUsers, err := s.transactionUC.GetTopTransactingUsers(c.Request.Context())
	if err != nil {
		log.Error().Err(err).Msg("TopUsers: failed to get top transacting users")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get list of top users",
		})
		return
	}
	c.JSON(http.StatusOK, topUsers)
}

func (s *HttpServer) TopTransactionsPerUser(c *gin.Context) {
	requestId := c.GetString("x-request-id")

	userToken := c.GetHeader("Authorization")
	if len(userToken) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "header 'Authorization' is empty",
		})
		return
	}

	username, err := s.userUC.GetUserNameFromToken(c.Request.Context(), userToken)
	if stacktrace.RootCause(err) == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	if err != nil {
		log.Error().Err(err).Ctx(c.Request.Context()).Str("requestId", requestId).Msg("Failed to get user by token")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get user information",
		})
		return
	}

	topTransactions, err := s.transactionUC.GetUserTopTransactions(c.Request.Context(), username)
	if err != nil {
		log.Error().Err(err).Msg("TopUsers: failed to get top user transaction")
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "failed to get list of top users",
		})
		return
	}
	c.JSON(http.StatusOK, topTransactions)
}
