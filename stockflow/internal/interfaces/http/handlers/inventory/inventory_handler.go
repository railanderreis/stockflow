package inventory

import (
	"github.com/gofiber/fiber/v2"
	appInventory "github.com/railanderreis/stockflow/stockflow/internal/application/inventory"
)

type InventoryHandler struct {
	inventoryUseCase *appInventory.InventoryUseCase
	warehouseUseCase *appInventory.WarehouseUseCase
}

func NewInventoryHandler(
	inventoryUseCase *appInventory.InventoryUseCase,
	warehouseUseCase *appInventory.WarehouseUseCase,
) *InventoryHandler {
	return &InventoryHandler{
		inventoryUseCase: inventoryUseCase,
		warehouseUseCase: warehouseUseCase,
	}
}

// GET /api/v1/warehouses
func (h *InventoryHandler) ListWarehouses(c *fiber.Ctx) error {
	warehouses, err := h.warehouseUseCase.ListAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(warehouses)
}

// POST /api/v1/warehouses
func (h *InventoryHandler) CreateWarehouse(c *fiber.Ctx) error {
	var input appInventory.CreateWarehouseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}

	warehouse, err := h.warehouseUseCase.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(warehouse)
}

// GET /api/v1/inventory/items
func (h *InventoryHandler) GetStockBalance(c *fiber.Ctx) error {
	warehouseID := c.Query("warehouse_id", "")
	productID := c.Query("product_id", "")

	balances, err := h.inventoryUseCase.GetBalance(c.Context(), warehouseID, productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(balances)
}

// POST /api/v1/inventory/adjustments
func (h *InventoryHandler) CreateAdjustment(c *fiber.Ctx) error {
	var input appInventory.CreateAdjustmentInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}
	input.CreatedByID = c.Locals("user_id").(string)

	adjustment, err := h.inventoryUseCase.CreateAdjustment(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(adjustment)
}

// POST /api/v1/inventory/transfers
func (h *InventoryHandler) TransferStock(c *fiber.Ctx) error {
	var input appInventory.TransferStockInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}
	input.RequestedByID = c.Locals("user_id").(string)

	transfer, err := h.inventoryUseCase.TransferStock(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(transfer)
}
