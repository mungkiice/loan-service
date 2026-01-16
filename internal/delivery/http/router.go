package http

import (
	"github.com/gin-gonic/gin"
	"github.com/mungkiice/-loan-service/internal/usecase"
)

func SetupRouter(handler *Handler, authHandler *AuthHandler, authUseCase *usecase.AuthUseCase) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/auth/signin", authHandler.SignIn)

		api.POST("/loans", handler.CreateLoan)
		api.GET("/loans", handler.GetLoans)
		api.GET("/loans/:id", handler.GetLoan)
	}

	protected := api.Group("")
	protected.Use(AuthMiddleware(authUseCase))
	{
		employeeRoutes := protected.Group("")
		employeeRoutes.Use(RequireUserType("employee"))
		{
			employeeRoutes.POST("/loans/:id/approve", RequireRole("field_validator"), handler.ApproveLoan)
			employeeRoutes.POST("/loans/:id/disburse", RequireRole("field_officer"), handler.DisburseLoan)
		}

		investorRoutes := protected.Group("")
		investorRoutes.Use(RequireUserType("investor"))
		{
			investorRoutes.POST("/loans/:id/invest", handler.Invest)
		}
	}

	return router
}
