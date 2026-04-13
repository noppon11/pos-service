package routes

import (
	"pos-service/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, posHandler *handler.PosHandler) {
	r.GET("/health", posHandler.GetHealth)
	r.GET("/api/v1/tenants/:tenant_id/health", posHandler.GetHealthByTenantID)
	r.GET("/api/v1/tenants/:tenant_id/branches", posHandler.GetBranchesByTenantID)
	r.GET("/api/v1/tenants/:tenant_id/branches/:branch_id", posHandler.GetByTenantIDAndBranchID)

}