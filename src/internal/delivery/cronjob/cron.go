package cronjob

import (
	"context"

	"github.com/idzharbae/digital-wallet/src/internal/repository"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog/log"
)

type Cron struct {
	userTransactionRepo repository.UserTransactionRepository
	cron                *cron.Cron
}

func NewCron(cron *cron.Cron, userTransactionRepo repository.UserTransactionRepository) *Cron {
	return &Cron{
		userTransactionRepo: userTransactionRepo,
		cron:                cron,
	}
}

func (c *Cron) Start() {
	c.cron.AddFunc("* * * * *", func() {
		_, err := c.userTransactionRepo.RefreshTopTransactingUsers(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("RefreshTopTransactingUsers Cron: error occured during data refresh")
		}
	})
	c.cron.Start()
}
