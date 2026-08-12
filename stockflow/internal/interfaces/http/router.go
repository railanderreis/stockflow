package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/auth"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/customer"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/inventory"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/product"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/purchase"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/reports"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/handlers/sales"
	"github.com/railanderreis/stockflow/stockflow/internal/interfaces/http/middleware"
)

type RouterDependencies struct {
	AuthHandler      *auth.AuthHandler
	ProductHandler   *product.ProductHandler
	InventoryHandler *inventory.InventoryHandler
	PurchaseHandler  *purchase.PurchaseHandler
	CustomerHandler  *customer.CustomerHandler
	SalesHandler     *sales.SalesHandler
	ReportsHandler   *reports.ReportsHandler
	AuthMiddleware   *middleware.AuthMiddleware
	AuditMiddleware  fiber.Handler
}

func SetupRoutes(app *fiber.App, deps RouterDependencies) {
	// Health check público
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "healthy",
			"service": "stockflow-api",
		})
	})

	// Prefix /api/v1
	v1 := app.Group("/api/v1")

	// -----------------------------------------------------------------
	// 1. AUTH & USER MANAGEMENT (Público / Autenticado)
	// -----------------------------------------------------------------
	authGroup := v1.Group("/auth")
	authGroup.Post("/login", deps.AuthHandler.Login)
	authGroup.Post("/refresh", deps.AuthHandler.RefreshToken)
	authGroup.Post("/register", deps.AuthMiddleware.RequirePermission("user:create"), deps.AuthHandler.Register)

	// Rotas Autenticadas Protegidas por JWT
	protected := v1.Group("", deps.AuthMiddleware.Authenticate())
	protected.Use(deps.AuditMiddleware) // Middleware global de Auditoria

	// -----------------------------------------------------------------
	// 2. PRODUCTS & CATEGORIES
	// -----------------------------------------------------------------
	categories := protected.Group("/categories")
	categories.Get("/", deps.ProductHandler.ListCategories)
	categories.Post("/", deps.AuthMiddleware.RequirePermission("category:create"), deps.ProductHandler.CreateCategory)

	products := protected.Group("/products")
	products.Get("/", deps.ProductHandler.ListProducts)
	products.Get("/:id", deps.ProductHandler.GetProductByID)
	products.Post("/", deps.AuthMiddleware.RequirePermission("product:create"), deps.ProductHandler.CreateProduct)
	products.Put("/:id", deps.AuthMiddleware.RequirePermission("product:update"), deps.ProductHandler.UpdateProduct)
	products.Delete("/:id", deps.AuthMiddleware.RequirePermission("product:delete"), deps.ProductHandler.DeleteProduct)

	// -----------------------------------------------------------------
	// 3. WAREHOUSES & INVENTORY MOVEMENTS
	// -----------------------------------------------------------------
	warehouses := protected.Group("/warehouses")
	warehouses.Get("/", deps.InventoryHandler.ListWarehouses)
	warehouses.Post("/", deps.AuthMiddleware.RequirePermission("warehouse:create"), deps.InventoryHandler.CreateWarehouse)

	inv := protected.Group("/inventory")
	inv.Get("/items", deps.InventoryHandler.GetStockBalance)
	inv.Post("/adjustments", deps.AuthMiddleware.RequirePermission("inventory:adjust"), deps.InventoryHandler.CreateAdjustment)
	inv.Post("/transfers", deps.AuthMiddleware.RequirePermission("inventory:transfer"), deps.InventoryHandler.TransferStock)

	// -----------------------------------------------------------------
	// 4. PROCUREMENT & PURCHASES
	// -----------------------------------------------------------------
	suppliers := protected.Group("/suppliers")
	suppliers.Get("/", deps.PurchaseHandler.ListSuppliers)
	suppliers.Post("/", deps.AuthMiddleware.RequirePermission("supplier:create"), deps.PurchaseHandler.CreateSupplier)

	purchases := protected.Group("/purchases/orders")
	purchases.Get("/", deps.PurchaseHandler.ListOrders)
	purchases.Get("/:id", deps.PurchaseHandler.GetOrderByID)
	purchases.Post("/", deps.AuthMiddleware.RequirePermission("purchase:create"), deps.PurchaseHandler.CreateOrder)
	purchases.Post("/:id/approve", deps.AuthMiddleware.RequirePermission("purchase:approve"), deps.PurchaseHandler.ApproveOrder)
	purchases.Post("/:id/receive", deps.AuthMiddleware.RequirePermission("purchase:receive"), deps.PurchaseHandler.ReceiveOrder)

	// -----------------------------------------------------------------
	// 5. CUSTOMERS & COMMERCIAL SALES
	// -----------------------------------------------------------------
	customers := protected.Group("/customers")
	customers.Get("/", deps.CustomerHandler.ListCustomers)
	customers.Post("/", deps.AuthMiddleware.RequirePermission("customer:create"), deps.CustomerHandler.CreateCustomer)

	salesOrders := protected.Group("/sales/orders")
	salesOrders.Get("/", deps.SalesHandler.ListOrders)
	salesOrders.Get("/:id", deps.SalesHandler.GetOrderByID)
	salesOrders.Post("/", deps.AuthMiddleware.RequirePermission("sales:create"), deps.SalesHandler.CreateOrder)
	salesOrders.Post("/:id/confirm", deps.AuthMiddleware.RequirePermission("sales:confirm"), deps.SalesHandler.ConfirmOrder)
	salesOrders.Post("/:id/ship", deps.AuthMiddleware.RequirePermission("sales:ship"), deps.SalesHandler.ShipOrder)
	salesOrders.Post("/:id/cancel", deps.AuthMiddleware.RequirePermission("sales:cancel"), deps.SalesHandler.CancelOrder)

	// -----------------------------------------------------------------
	// 6. REPORTS & AUDIT LOGS
	// -----------------------------------------------------------------
	rep := protected.Group("/reports")
	rep.Get("/abc", deps.AuthMiddleware.RequirePermission("reports:view"), deps.ReportsHandler.GetABCAnalysis)
	rep.Get("/valuation", deps.AuthMiddleware.RequirePermission("reports:view"), deps.ReportsHandler.GetValuation)
	rep.Get("/reorder-alerts", deps.AuthMiddleware.RequirePermission("reports:view"), deps.ReportsHandler.GetReorderAlerts)

	protected.Get("/audit-logs", deps.AuthMiddleware.RequirePermission("audit:view"), deps.ReportsHandler.GetAuditLogs)
}
