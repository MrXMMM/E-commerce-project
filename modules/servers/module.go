package servers

import (
	middlewareshandlers "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresUsecases"
	"github.com/MrXMMM/E-commerce-Project/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewareshandlers.IMiddlewaresHandler
}

func InitModule(router fiber.Router, server *server, mid middlewareshandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router: router,
		server: server,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewareshandlers.IMiddlewaresHandler {
	repository := middlewaresrepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresusecases.MiddlewaresUsecase(repository)
	return middlewareshandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}
