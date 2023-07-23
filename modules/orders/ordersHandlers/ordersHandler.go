package ordershandlers

import (
	"fmt"
	"strings"
	"time"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
	"github.com/MrXMMM/E-commerce-Project/modules/orders"
	ordersusecases "github.com/MrXMMM/E-commerce-Project/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type orderHandlerErrCode string

const (
	findOneOrderErr orderHandlerErrCode = "order-001"
	findOrderErr    orderHandlerErrCode = "order-002"
	insertOrderErr  orderHandlerErrCode = "order-003"
	updateOrderErr  orderHandlerErrCode = "order-004"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	cfg           config.IConfig
	ordersUsecase ordersusecases.IOrdersUsecase
}

func OrdersHandler(cfg config.IConfig, ordersUsecase ordersusecases.IOrdersUsecase) IOrdersHandler {
	return &ordersHandler{
		cfg:           cfg,
		ordersUsecase: ordersUsecase,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {

	orderId := strings.Trim(c.Params("order_id"), " ")

	order, err := h.ordersUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *ordersHandler) FindOrder(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq:       &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},
	}
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 5 {
		req.Limit = 5
	}

	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}

	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}

	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(findOneOrderErr),
				"start date is invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(findOneOrderErr),
				"end date is invalid",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	// Usecase
	orders := h.ordersUsecase.FindOrder(req)

	return entities.NewResponse(c).Success(fiber.StatusOK, orders).Res()
}

func (h *ordersHandler) InsertOrder(c *fiber.Ctx) error {

	userId := c.Locals("userId").(string)

	req := &orders.Order{
		Product: make([]*orders.ProductsOrder, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	fmt.Println(req)
	if len(req.Product) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			"products are empty",
		).Res()
	}

	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.ordersUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, order).Res()
}

func (h *ordersHandler) UpdateOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	req := new(orders.Order)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}

	req.Id = orderId

	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled",
	}

	fmt.Println(c.Locals("userRoleId").(int))
	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else {
		req.Status = statusMap["canceled"]
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}
		if req.TransferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(updateOrderErr),
					err.Error(),
				).Res()
			}

			now := time.Now().In(loc)

			// YYYY-MM-DD HH:MM:SS
			// 2006-01-02 15:04:05
			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")
		}
	}

	order, err := h.ordersUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
