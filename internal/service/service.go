package service

import (
	"utopia/internal/config"

	"github.com/gin-gonic/gin"
)

type Service struct {
	Config *config.ServiceConfig
	Router *gin.Engine
}

func NewService(config *config.ServiceConfig) *Service {
	return &Service{
		Config: config,
	}
}

func (s *Service) Start() error {
	return nil
}

func (s *Service) Stop() error {
	return nil
}
