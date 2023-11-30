package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/controllers/admin"
	"github.com/myanmarmarathon/mkitchen-distribution-backend/middlewares"
)

func SetupRoutes(r *gin.Engine) {
    
	r.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})
	authRouter := r.Group("/api/v1")
	authRouter.POST("/login", admin.Login)

	protectedRouter := r.Group("/api/v1")
	protectedRouter.Use(middlewares.JwtAuthMiddleware())

	protectedRouter.GET("/profile", admin.CurrentUser)
	protectedRouter.POST("/logout", admin.Logout)

	protectedRouter.GET("/users", admin.GetAllUsers)
	protectedRouter.POST("/users", admin.CreateUser)
	protectedRouter.PATCH("/users/:id", admin.UpdateUser)
	protectedRouter.DELETE("/users/:id", admin.DeleteUser)
	protectedRouter.GET("/users/:id", admin.GetUser)

	protectedRouter.GET("/product_categories", admin.GetAllProductCategories)
	protectedRouter.POST("/product_categories", admin.CreateProductCategory)
	protectedRouter.PATCH("/product_categories/:id", admin.UpdateProductCategory)
	protectedRouter.DELETE("/product_categories/:id", admin.DeleteProductCategory)
	protectedRouter.GET("/product_categories/:id", admin.GetProductCategory)

	protectedRouter.GET("/suppliers", admin.GetAllSuppliers)
	protectedRouter.POST("/suppliers", admin.CreateSupplier)
	protectedRouter.PATCH("/suppliers/:id", admin.UpdateSupplier)
	protectedRouter.DELETE("/suppliers/:id", admin.DeleteSupplier)
	protectedRouter.GET("/suppliers/:id", admin.GetSupplier)

	protectedRouter.GET("/products", admin.GetAllProducts)
	protectedRouter.POST("/products", admin.CreateProduct)
	protectedRouter.PATCH("/products/:id", admin.UpdateProduct)
	protectedRouter.DELETE("/products/:id", admin.DeleteProduct)
	protectedRouter.GET("/products/:id", admin.GetProduct)

	protectedRouter.POST("/upload_image", admin.UploadImage)
	protectedRouter.DELETE("/delete_image/:id", admin.DeleteImage)

	protectedRouter.GET("/purchase_orders", admin.GetAllPurchaseOrders)
	protectedRouter.POST("/purchase_orders", admin.CreatePurchaseOrder)
	protectedRouter.PATCH("/purchase_orders/:id", admin.UpdatePurchaseOrder)
	protectedRouter.DELETE("/purchase_orders/:id", admin.DeletePurchaseOrder)
	protectedRouter.GET("/purchase_orders/:id", admin.GetPurchaseOrder)

	protectedRouter.POST("/purchase_orders/:id/receive", admin.ReceivePurchaseOrder)
}