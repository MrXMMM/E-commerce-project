package servers

import (
	appinfohandlers "github.com/MrXMMM/E-commerce-Project/modules/appinfo/appinfoHandlers"
	appinforepositories "github.com/MrXMMM/E-commerce-Project/modules/appinfo/appinfoRepositories"
	appinfousecases "github.com/MrXMMM/E-commerce-Project/modules/appinfo/appinfoUsecases"
	fileshandlers "github.com/MrXMMM/E-commerce-Project/modules/files/filesHandlers"
	filesusecases "github.com/MrXMMM/E-commerce-Project/modules/files/filesUsecases"
	middlewareshandlers "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresHandlers"
	middlewaresrepositories "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresRepositories"
	middlewaresusecases "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresUsecases"
	"github.com/MrXMMM/E-commerce-Project/modules/monitor/monitorHandlers"
	productshandlers "github.com/MrXMMM/E-commerce-Project/modules/products/productsHandlers"
	productsrepositories "github.com/MrXMMM/E-commerce-Project/modules/products/productsRepositories"
	productsusecases "github.com/MrXMMM/E-commerce-Project/modules/products/productsUsecases"
	usershandlers "github.com/MrXMMM/E-commerce-Project/modules/users/usersHandlers"
	usersrepositories "github.com/MrXMMM/E-commerce-Project/modules/users/usersRepositories"
	usersusecases "github.com/MrXMMM/E-commerce-Project/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
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

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/admin/signup", m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:userid", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

	// Initial admin 1 person in db (SQL script)
	// Generate Admin Key
	// Everytime using insert admin. It need to send the admin key by middleware
}

func (m *moduleFactory) AppinfoModule() {
	repository := appinforepositories.AppinfoRepository(m.server.db)
	usecase := appinfousecases.AppinfoUsecase(repository)
	handler := appinfohandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)

	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)

	router.Delete("/:category_id/category", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCategory)
}

func (m *moduleFactory) FilesModule() {
	usecase := filesusecases.FileUsecase(m.server.cfg)
	handler := fileshandlers.FilesHandler(m.server.cfg, usecase)

	router := m.router.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)
}

func (m *moduleFactory) ProductsModule() {
	fileUsecase := filesusecases.FileUsecase(m.server.cfg)

	repository := productsrepositories.ProductsRepository(m.server.db, m.server.cfg, fileUsecase)
	usecase := productsusecases.ProductsUsecase(repository)
	handler := productshandlers.ProductHandler(m.server.cfg, usecase, fileUsecase)

	router := m.router.Group("/products")

	router.Get("/", m.mid.ApiKeyAuth(), handler.FindProduct)
	router.Get("/:product_id", m.mid.ApiKeyAuth(), handler.FindOneProduct)

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddProduct)

	router.Patch("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateProduct)

	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteProduct)

}
