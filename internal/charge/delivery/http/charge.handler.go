package handler

import (
	"herbie-charge-api/internal/charge/dto"
	"herbie-charge-api/internal/charge/service"

	"github.com/gofiber/fiber/v2"
)

type ChargeHandler struct {
	service *service.ChargeService
}

func NewChargeHandler(
	service *service.ChargeService,
) *ChargeHandler {
	return &ChargeHandler{
		service: service,
	}
}

func (h *ChargeHandler) Charge(
	c *fiber.Ctx,
) error {
	var req dto.ChargeRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			dto.ChargeResponse{
				Success: false,
				Message: "invalid request",
			},
		)
	}

	result, err :=
		h.service.Charge(req)

	if err != nil {
		return c.Status(
			fiber.StatusBadRequest,
		).JSON(
			dto.ChargeResponse{
				Success: false,
				Message: err.Error(),
			},
		)
	}

	return c.JSON(
		dto.ChargeResponse{
			Success: true,
			Message: "เติมเครดิตสำเร็จ",
			Data:    result,
		},
	)
}
