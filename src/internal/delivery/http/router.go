package http

func (s *HttpServer) SetupRouters() {
	v1 := s.engine.Group("/v1")
	{
		v1.GET("/ping", s.Ping)
		v1.GET("/messages", s.ListMessages)
		v1.POST("/publish/example", s.SendMessage)
		v1.POST("/create_user", s.RegisterUser)
		v1.GET("/balance_read", s.BalanceRead)
	}
}
