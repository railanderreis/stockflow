package reports

import (
	"github.com/gofiber/fiber/v2"
	appReports "github.com/railanderreis/stockflow/stockflow/internal/application/reports"
)

type ReportsHandler struct {
	reportsUseCase *appReports.ReportsUseCase
}

func NewReportsHandler(reportsUseCase *appReports.ReportsUseCase) *ReportsHandler {
	return &ReportsHandler{
		reportsUseCase: reportsUseCase,
	}
}

// GET /api/v1/reports/abc
func (h *ReportsHandler) GetABCAnalysis(c *fiber.Ctx) error {
	days := c.QueryInt("days", 30)

	analysis, err := h.reportsUseCase.GetABCAnalysis(c.Context(), days)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(analysis)
}

// GET /api/v1/reports/valuation
func (h *ReportsHandler) GetValuation(c *fiber.Ctx) error {
	warehouseID := c.Query("warehouse_id", "")

	valuation, err := h.reportsUseCase.GetValuation(c.Context(), warehouseID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(valuation)
}

// GET /api/v1/reports/reorder-alerts
func (h *ReportsHandler) GetReorderAlerts(c *fiber.Ctx) error {
	alerts, err := h.reportsUseCase.GetReorderAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(alerts)
}

// GET /api/v1/audit-logs
func (h *ReportsHandler) GetAuditLogs(c *fiber.Ctx) error {
	entity := c.Query("entity", "")
	userID := c.Query("user_id", "")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 50)

	logs, err := h.reportsUseCase.GetAuditLogs(c.Context(), entity, userID, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(logs)
}
