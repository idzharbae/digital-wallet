package http

func (s *HttpServer) SetupRouters() {
	v1 := s.engine.Group("/v1")
	{
		v1.GET("/ping", s.Ping)
		v1.GET("/messages", s.ListMessages)
		v1.POST("/publish/example", s.SendMessage)
		v1.POST("/create_user", s.RegisterUser)
		v1.GET("/balance_read", s.BalanceRead)
		v1.POST("/balance_topup", s.BalanceTopUp)
		v1.POST("/transfer", s.Transfer)
		v1.GET("/top_users", s.TopUsers)
		v1.GET("/top_transactions_per_user", s.TopTransactionsPerUser)
	}
}
