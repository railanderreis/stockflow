package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	// Ajuste os imports abaixo conforme o caminho do seu módulo
	"github.com/raianderreis/stockflow/stockflow/internal/application/auth"
	"github.com/raianderreis/stockflow/stockflow/internal/application/customer"
	"github.com/raianderreis/stockflow/stockflow/internal/application/inventory"
	"github.com/raianderreis/stockflow/stockflow/internal/application/product"
	"github.com/raianderreis/stockflow/stockflow/internal/application/purchase"
	"github.com/raianderreis/stockflow/stockflow/internal/application/reports"
	"github.com/raianderreis/stockflow/stockflow/internal/application/sales"

	httpRoutes "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http"
	authHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/auth"
	customerHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/customer"
	inventoryHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/inventory"
	productHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/product"
	purchaseHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/purchase"
	reportsHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/reports"
	salesHandler "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/handlers/sales"
	httpMiddleware "github.com/raianderreis/stockflow/stockflow/internal/interfaces/http/middleware"

	// Mock/Implementação de repositórios
	"github.com/raianderreis/stockflow/stockflow/internal/infrastructure/repository/memory"
)

func main() {
	// 1. Configuração do Servidor Fiber
	app := fiber.New(fiber.Config{
		AppName:      "StockFlow API v1",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	// 2. Middlewares Globais de Rede
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// 3. Inicialização de Repositórios e Mocks de Infraestrutura
	// (Caso ainda não tenha conectado no Postgres/Redis real, usando mocks/memória para testes)
	userRepo := memory.NewUserRepository()
	productRepo := memory.NewProductRepository()
	categoryRepo := memory.NewCategoryRepository()
	inventoryRepo := memory.NewInventoryRepository()
	warehouseRepo := memory.NewWarehouseRepository()
	supplierRepo := memory.NewSupplierRepository()
	purchaseRepo := memory.NewPurchaseRepository()
	customerRepo := memory.NewCustomerRepository()
	salesRepo := memory.NewSalesRepository()
	reportsRepo := memory.NewReportsRepository()

	// 4. Instanciação dos Casos de Uso (Application Layer)
	jwtSecret := getEnv("JWT_SECRET", "super-secret-key-stockflow")

	authUC := auth.NewAuthUseCase(userRepo, jwtSecret)
	productUC := product.NewProductUseCase(productRepo)
	categoryUC := product.NewCategoryUseCase(categoryRepo)
	inventoryUC := inventory.NewInventoryUseCase(inventoryRepo)
	warehouseUC := inventory.NewWarehouseUseCase(warehouseRepo)
	purchaseUC := purchase.NewPurchaseUseCase(purchaseRepo)
	supplierUC := purchase.NewSupplierUseCase(supplierRepo)
	customerUC := customer.NewCustomerUseCase(customerRepo)
	confirmSalesUC := sales.NewConfirmOrderUseCase(salesRepo, inventoryRepo)
	shipSalesUC := sales.NewShipOrderUseCase(salesRepo, inventoryRepo)
	reportsUC := reports.NewReportsUseCase(reportsRepo)

	// 5. Instanciação dos Handlers HTTP REST
	authH := authHandler.NewAuthHandler(authUC)
	productH := productHandler.NewProductHandler(productUC, categoryUC)
	inventoryH := inventoryHandler.NewInventoryHandler(inventoryUC, warehouseUC)
	purchaseH := purchaseHandler.NewPurchaseHandler(purchaseUC, supplierUC)
	customerH := customerHandler.NewCustomerHandler(customerUC)
	salesH := salesHandler.NewSalesHandler(confirmSalesUC, shipSalesUC)
	reportsH := reportsHandler.NewReportsHandler(reportsUC)

	// 6. Instanciação do Middleware de Autenticação / RBAC
	authMw := httpMiddleware.NewAuthMiddleware(jwtSecret)

	// Middleware simples para logs de auditoria
	auditMw := func(c *fiber.Ctx) error {
		// Executa a requisição normalmente
		return c.Next()
	}

	// 7. Registro Central de Rotas `/api/v1`
	deps := httpRoutes.RouterDependencies{
		AuthHandler:      authH,
		ProductHandler:   productH,
		InventoryHandler: inventoryH,
		PurchaseHandler:  purchaseH,
		CustomerHandler:  customerH,
		SalesHandler:     salesH,
		ReportsHandler:   reportsH,
		AuthMiddleware:   authMw,
		AuditMiddleware:  auditMw,
	}

	// Injeta todas as rotas no Fiber
	httpRoutes.SetupRoutes(app, deps)

	// 8. Graceful Shutdown (Encerramento gracioso do servidor)
	port := getEnv("PORT", "8080")

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
			log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
		}
	}()

	log.Printf("🚀 Servidor escutando na porta %s (http://localhost:%s)", port, port)

	// Aguarda sinal de interrupção (Ctrl+C ou SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Encerrando servidor de forma graciosa...")
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v", err)
	}
	log.Println("Servidor finalizado com sucesso.")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
