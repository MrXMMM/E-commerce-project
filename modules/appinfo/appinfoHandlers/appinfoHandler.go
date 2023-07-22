package appinfohandlers

import (
	"strconv"
	"strings"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/appinfo"
	appinfousecases "github.com/MrXMMM/E-commerce-Project/modules/appinfo/appinfoUsecases"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
	"github.com/MrXMMM/E-commerce-Project/pkg/auth"
	"github.com/gofiber/fiber/v2"
)

type appinfoHandlerErrCode string

const (
	generateApiKeyErr appinfoHandlerErrCode = "appinfo-001"
	findCategoryErr   appinfoHandlerErrCode = "appinfo-002"
	addCategoryErr    appinfoHandlerErrCode = "appinfo-003"
	removeCategoryErr appinfoHandlerErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	AddCategory(c *fiber.Ctx) error
	RemoveCategory(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsecases appinfousecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecases appinfousecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsecases: appinfoUsecases,
	}
}

func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := auth.NewAuth(
		auth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
}

func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CatagoryFilter)
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsecases.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, category).Res()
}

func (h *appinfoHandler) AddCategory(c *fiber.Ctx) error {

	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addCategoryErr),
			"categories request are empty",
		).Res()
	}

	if err := h.appinfoUsecases.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(addCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) RemoveCategory(c *fiber.Ctx) error {

	categoryId := strings.Trim(c.Params("category_id"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			"id type is invalid",
		).Res()
	}

	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			"id must be more than 0",
		).Res()
	}

	if err := h.appinfoUsecases.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(removeCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		&struct {
			CategoryId int `่json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}
