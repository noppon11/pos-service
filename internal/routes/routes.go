package routes

import (
	"pos-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, posHandler *handler.PosHandler) {
	r.GET("/health", posHandler.GetHealth)
	r.GET("/readiness", posHandler.Readiness)
	api := r.Group("/api/v1")
	{
		api.GET("/tenants/:tenant_id/health", posHandler.GetHealthByTenantID)
		api.GET("/tenants/:tenant_id/branches", posHandler.GetBranchesByTenantID)
		api.GET("/tenants/:tenant_id/branches/:branch_id", posHandler.GetByTenantIDAndBranchID)
	}
}