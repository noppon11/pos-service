package routes

import (
	"pos-service/internal/handler"
	"pos-service/internal/middleware"
	"pos-service/internal/service"
	"pos-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	posHandler *handler.PosHandler,
	authHandler *handler.AuthHandler,
	authService *service.AuthService,
	jwtManager *utils.JWTManager,
) {
	r.GET("/health", posHandler.GetHealth)
	r.GET("/readiness", posHandler.Readiness)

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/me",
				middleware.RequireAuth(jwtManager),
				authHandler.Me,
			)
		}

		protectedProducts := api.Group("/tenants/:tenantID/branches/:branchID")
		protectedProducts.Use(middleware.RequireAuth(jwtManager))
		protectedProducts.Use(middleware.RequireTenantBranchAccess(authService))
		{
			protectedProducts.GET("/products", posHandler.GetAllProducts)
			protectedProducts.POST("/products", posHandler.CreateProduct)
			protectedProducts.GET("/products/:productID", posHandler.GetProductByID)
			protectedProducts.PUT("/products/:productID", posHandler.UpdateProduct)
			protectedProducts.DELETE("/products/:productID", posHandler.DeleteProduct)
		}
	}
}