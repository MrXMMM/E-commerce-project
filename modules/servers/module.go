package servers

import (
	middlewareshandlers "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresUsecases"
	"github.com/MrXMMM/E-commerce-Project/modules/monitor/monitorHandlers"
	usershandlers "github.com/MrXMMM/E-commerce-Project/modules/users/usersHandlers"
	usersrepositories "github.com/MrXMMM/E-commerce-Project/modules/users/usersRepositories"
	usersusecases "github.com/MrXMMM/E-commerce-Project/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
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

func (m *moduleFactory) UsersModule() {
	repository := usersrepositories.UserRepository(m.server.db)
	usecase := usersusecases.UsersUseCase(m.server.cfg, repository)
	handler := usershandlers.UsersHandler(m.server.cfg, usecase)

	// v1/users

	router := m.router.Group("/users")

	router.Post("/signup", handler.SignUpCustomer)
	router.Post("/signin", handler.SignIn)
	router.Post("/refresh", handler.RefreshPassport)
	router.Post("/signout", handler.SignOut)
	router.Post("/admin/signup", m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:userid", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	// Initial admin 1 person in db (SQL script)
	// Generate Admin Key
	// Everytime using insert admin. It need to send the admin key by middleware
}
