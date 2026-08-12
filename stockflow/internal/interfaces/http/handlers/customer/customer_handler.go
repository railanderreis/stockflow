package customer

import (
	"github.com/gofiber/fiber/v2"
	appCustomer "github.com/railanderreis/stockflow/stockflow/internal/application/customer"
)

type CustomerHandler struct {
	customerUseCase *appCustomer.CustomerUseCase
}

func NewCustomerHandler(customerUseCase *appCustomer.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{
		customerUseCase: customerUseCase,
	}
}

// GET /api/v1/customers
func (h *CustomerHandler) ListCustomers(c *fiber.Ctx) error {
	search := c.Query("search", "")
	customers, err := h.customerUseCase.List(c.Context(), search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(customers)
}

// POST /api/v1/customers
func (h *CustomerHandler) CreateCustomer(c *fiber.Ctx) error {
	var input appCustomer.CreateCustomerInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}

	customer, err := h.customerUseCase.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}
