package sales

import (
	"github.com/gofiber/fiber/v2"
	appSales "github.com/railanderreis/stockflow/stockflow/internal/application/sales"
)

type SalesHandler struct {
	confirmOrderUseCase *appSales.ConfirmOrderUseCase
	shipOrderUseCase    *appSales.ShipOrderUseCase
}

func NewSalesHandler(
	confirmOrderUseCase *appSales.ConfirmOrderUseCase,
	shipOrderUseCase *appSales.ShipOrderUseCase,
) *SalesHandler {
	return &SalesHandler{
		confirmOrderUseCase: confirmOrderUseCase,
		shipOrderUseCase:    shipOrderUseCase,
	}
}

// POST /api/v1/sales/orders/:id/confirm
func (h *SalesHandler) ConfirmOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "order ID is required"})
	}

	if err := h.confirmOrderUseCase.Execute(c.Context(), orderID); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sales order confirmed and stock successfully reserved",
	})
}

// POST /api/v1/sales/orders/:id/ship
func (h *SalesHandler) ShipOrder(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var input appSales.ShipOrderInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request payload"})
	}
	input.SalesOrderID = orderID
	input.DispatchedByID = c.Locals("user_id").(string)

	shipment, err := h.shipOrderUseCase.Execute(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(shipment)
}

func (h *SalesHandler) ListOrders(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *SalesHandler) GetOrderByID(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *SalesHandler) CreateOrder(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}

func (h *SalesHandler) CancelOrder(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotImplemented)
}
