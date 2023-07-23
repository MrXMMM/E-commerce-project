package middlewareshandlers

import (
	"strings"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
	middlewaresusecases "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresUsecases"
	"github.com/MrXMMM/E-commerce-Project/pkg/auth"
	"github.com/MrXMMM/E-commerce-Project/pkg/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlerErrCode string

const (
	routerCheckErr middlewareHandlerErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlerErrCode = "middleware-002"
	paramsCheckErr middlewareHandlerErrCode = "middleware-003"
	authorizeErr   middlewareHandlerErrCode = "middleware-004"
	apiKeyErr      middlewareHandlerErrCode = "middlware-005"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRole ...int) fiber.Handler
	ApiKeyAuth() fiber.Handler
}

type middlewaresHandler struct {
	cfg                config.IConfig
	middlewaresUseCase middlewaresusecases.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, middlewaresUseCase middlewaresusecases.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:                cfg,
		middlewaresUseCase: middlewaresUseCase,
	}
}

func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		if !h.middlewaresUseCase.FindAccessToken(claims.Id, token) {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"no permission to access",
			).Res()
		}

		// Set UserId
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("userid") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"You not have an access to view this",
			).Res()
		}
		return c.Next()
	}
}

func (h *middlewaresHandler) Authorize(expectRole ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("user_roleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user_id is not int type",
			).Res()
		}

		roles, err := h.middlewaresUseCase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRole {
			sum += v
		}

		expectedValueBinary := util.BinaryConverter(sum, len(roles))
		userValueBinary := util.BinaryConverter(userRoleId, len(roles))

		for i := range userValueBinary {
			if userValueBinary[i]&expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"this role is not permission for this api",
		).Res()

	}
}

func (h *middlewaresHandler) ApiKeyAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		key := c.Get("X-Api-Key")
		if _, err := auth.ParseApiKey(h.cfg.Jwt(), key); err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(apiKeyErr),
				"apikey is invalid or required",
			).Res()
		}
		return c.Next()
	}
}
