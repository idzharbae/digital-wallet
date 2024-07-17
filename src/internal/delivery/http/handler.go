package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/idzharbae/digital-wallet/src/internal/constants"
	"github.com/idzharbae/digital-wallet/src/internal/delivery/http/dto"
	"github.com/idzharbae/digital-wallet/src/internal/entity"
	"github.com/idzharbae/digital-wallet/src/internal/utils"
	"github.com/palantir/stacktrace"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func (s *HttpServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (s *HttpServer) SendMessage(c *gin.Context) {
	var msg entity.Message

	request_id := c.GetString("x-request-id")

	// Bind request payload with our model
	if binderr := c.ShouldBindJSON(&msg); binderr != nil {

		log.Error().Err(binderr).Str("request_id", request_id).
			Msg("Error occurred while binding request data")

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": binderr.Error(),
		})
		return
	}

	s.rmqProducer.PublishMessage(constants.EXAMPLE_QUEUE, "text/plain", []byte(msg.Message))

	c.JSON(http.StatusOK, gin.H{
		"response": "Message received",
	})
}

func (s *HttpServer) ListMessages(c *gin.Context) {
	type msg struct {
		Id      int    `json:"id"`
		Message string `json:"message"`
	}
	messages := []msg{}

	val, err := s.redisClient.Get(c.Request.Context(), constants.EXAMPLE_REDIS_KEY).Result()
	if err == nil {
		log.Info().Msg("redis hit")
		json.Unmarshal([]byte(val), &messages)
		c.JSON(http.StatusOK, messages)
		return
	}
	if err != redis.Nil {
		c.AbortWithError(http.StatusInternalServerError, stacktrace.Propagate(err, "failed get data from redis"))
		return
	}

	rows, err := s.dbConn.Query(c.Request.Context(), `SELECT id, message FROM dummy_table`)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, stacktrace.Propagate(err, "failed to query to dummy_table"))
		return
	}

	defer rows.Close()
	for rows.Next() {
		var message msg
		err := rows.Scan(&message.Id, &message.Message)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, stacktrace.Propagate(err, "failed scan row"))
			return
		}

		messages = append(messages, message)
	}

	marshaledMessages, _ := json.Marshal(messages)
	err = s.redisClient.Set(c.Request.Context(), constants.EXAMPLE_QUEUE, marshaledMessages, 0).Err()
	if err != nil {
		log.Error().Err(err).Msg("failed to set redis")
	}
	c.JSON(http.StatusOK, messages)
}

func (s *HttpServer) RegisterUser(c *gin.Context) {
	request_id := c.GetString("x-request-id")
	var request dto.RegisterUserRequest
	if binderr := c.ShouldBindJSON(&request); binderr != nil {
		log.Error().Err(binderr).Str("request_id", request_id).
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
			"message": "username must start with alphabet and must be alphanumeric",
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
