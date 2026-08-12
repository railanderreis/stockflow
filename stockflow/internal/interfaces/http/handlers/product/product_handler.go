package product

import (
	"github.com/gofiber/fiber/v2"
	appProduct "github.com/railanderreis/stockflow/stockflow/internal/application/product"
)

type ProductHandler struct {
	productUseCase  *appProduct.ProductUseCase
	categoryUseCase *appProduct.CategoryUseCase
}

func NewProductHandler(
	productUseCase *appProduct.ProductUseCase,
	categoryUseCase *appProduct.CategoryUseCase,
) *ProductHandler {
	return &ProductHandler{
		productUseCase:  productUseCase,
		categoryUseCase: categoryUseCase,
	}
}

// GET /api/v1/products
func (h *ProductHandler) ListProducts(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)
	search := c.Query("search", "")

	result, err := h.productUseCase.List(c.Context(), page, limit, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// GET /api/v1/products/:id
func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.productUseCase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "produto não encontrado"})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// POST /api/v1/products
func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var input appProduct.CreateProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}

	product, err := h.productUseCase.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var input appProduct.UpdateProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}

	product, err := h.productUseCase.Update(c.Context(), id, input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.productUseCase.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GET /api/v1/categories
func (h *ProductHandler) ListCategories(c *fiber.Ctx) error {
	categories, err := h.categoryUseCase.ListAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(categories)
}

// POST /api/v1/categories
func (h *ProductHandler) CreateCategory(c *fiber.Ctx) error {
	var input appProduct.CreateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "payload inválido"})
	}

	category, err := h.categoryUseCase.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(category)
}
