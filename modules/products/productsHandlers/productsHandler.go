package productshandlers

import (
	"fmt"
	"strings"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/appinfo"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
	"github.com/MrXMMM/E-commerce-Project/modules/files"
	filesusecases "github.com/MrXMMM/E-commerce-Project/modules/files/filesUsecases"
	"github.com/MrXMMM/E-commerce-Project/modules/products"
	productsusecases "github.com/MrXMMM/E-commerce-Project/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)

type productsHandlersErrCode string

const (
	findOneProductErr productsHandlersErrCode = "product-001"
	findProductErr    productsHandlersErrCode = "product-002"
	insertProductErr  productsHandlersErrCode = "product-003"
	updateProductErr  productsHandlersErrCode = "product-004"
	deleteProductErr  productsHandlersErrCode = "product-005"
)

type IProductsHndler interface {
	FindOneProduct(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
}

type productHandler struct {
	cfg            config.IConfig
	productUsecase productsusecases.IProductsUsecase
	filesUsecase   filesusecases.IFilesUsecase
}

func ProductHandler(cfg config.IConfig, productUsecase productsusecases.IProductsUsecase, filesUsecase filesusecases.IFilesUsecase) IProductsHndler {
	return &productHandler{
		cfg:            cfg,
		productUsecase: productUsecase,
		filesUsecase:   filesUsecase,
	}
}

func (h *productHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productHandler) FindProduct(c *fiber.Ctx) error {
	req := &products.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findProductErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products := h.productUsecase.FindProduct(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productHandler) AddProduct(c *fiber.Ctx) error {
	req := &products.Product{
		Category: &appinfo.Category{},
		Images:   make([]*entities.Image, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	if req.Category.Id <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			"category id is invalid",
		).Res()
	}

	product, err := h.productUsecase.AddProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productHandler) UpdateProduct(c *fiber.Ctx) error {

	productId := strings.Trim(c.Params("product_id"), " ")

	req := &products.Product{
		Images:   make([]*entities.Image, 0),
		Category: &appinfo.Category{},
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	req.Id = productId

	product, err := h.productUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, product).Res()
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("product_id"), " ")

	product, err := h.productUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	deleteFileReq := make([]*files.DeleteFileReq, 0)
	for _, p := range product.Images {
		deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
			Destination: fmt.Sprintf("images/product/%s", p.FileName),
		})
	}

	if err := h.filesUsecase.DeleteFileOnGCP(deleteFileReq); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	if err := h.productUsecase.DeleteProduct(productId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusNoContent, nil).Res()
}
