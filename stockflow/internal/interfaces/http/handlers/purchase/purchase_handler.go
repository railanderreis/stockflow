package purchase

import (
	"github.com/gofiber/fiber/v2"
	appPurchase "github.com/railanderreis/stockflow/stockflow/internal/application/purchase"
)

type PurchaseHandler struct {
	purchaseUseCase *appPurchase.PurchaseUseCase
	supplierUseCase *appPurchase.SupplierUseCase
}

func NewPurchaseHandler(
	purchaseUseCase *appPurchase.PurchaseUseCase,
	supplierUseCase *appPurchase.SupplierUseCase,
) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseUseCase: purchaseUseCase,
		supplierUseCase: supplierUseCase,
	}
}

// GET /api/v1/suppliers
func (h *PurchaseHandler) ListSuppliers(c *fiber.Ctx) error {
	suppliers, err := h.supplierUseCase.ListAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(suppliers)
}

// POST /api/v1/suppliers
func (h *PurchaseHandler) CreateSupplier(c *fiber.Ctx) error {
	var input appPurchase.CreateSupplierInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}

	supplier, err := h.supplierUseCase.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(supplier)
}

// GET /api/v1/purchases/orders
func (h *PurchaseHandler) ListOrders(c *fiber.Ctx) error {
	status := c.Query("status", "")
	orders, err := h.purchaseUseCase.ListOrders(c.Context(), status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

// GET /api/v1/purchases/orders/:id
func (h *PurchaseHandler) GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, err := h.purchaseUseCase.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "ordem de compra não encontrada"})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

// POST /api/v1/purchases/orders
func (h *PurchaseHandler) CreateOrder(c *fiber.Ctx) error {
	var input appPurchase.CreateOrderInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}
	input.CreatedByID = c.Locals("user_id").(string)

	order, err := h.purchaseUseCase.CreateOrder(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// POST /api/v1/purchases/orders/:id/approve
func (h *PurchaseHandler) ApproveOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	approvedBy := c.Locals("user_id").(string)

	if err := h.purchaseUseCase.ApproveOrder(c.Context(), id, approvedBy); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "ordem de compra aprovada com sucesso"})
}

// POST /api/v1/purchases/orders/:id/receive
func (h *PurchaseHandler) ReceiveOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var input appPurchase.ReceiveOrderInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}
	input.PurchaseOrderID = id
	input.ReceivedByID = c.Locals("user_id").(string)

	if err := h.purchaseUseCase.ReceiveOrder(c.Context(), input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "estoque recebido e custo médio atualizado"})
}
