package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go-ecommerce-app/internal/api/rest"
	"go-ecommerce-app/internal/repository"
	"go-ecommerce-app/internal/service"
	"log"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func SetupCatalogRoutes(rh *rest.RestHandler) {
	app := rh.App

	// Create an instance of user service and inject to handler
	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := CatalogHandler{
		svc: svc,
	}

	// Public Endpoints
	app.Get("/products")
	app.Get("/products/:id")
	app.Get("/categories")
	app.Get("/categories/:id")

	// Private Endpoints
	selRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller)
	// Categories
	selRoutes.Post("/categories", handler.CreateCategories)
	selRoutes.Patch("/categories/:id", handler.EditCategory)
	selRoutes.Delete("/categories/:id", handler.DeleteCategory)
	// Products
	selRoutes.Post("/products", handler.CreateProducts)
	selRoutes.Get("/products", handler.GetProducts)
	selRoutes.Get("/products/:id", handler.GetProduct)
	selRoutes.Patch("/products/:id", handler.EditProduct)
	selRoutes.Put("/products/:id", handler.UpdateProduct)
	selRoutes.Delete("/products/:id", handler.DeleteProduct)
}

// Categories

func (h CatalogHandler) CreateCategories(ctx *fiber.Ctx) error {
	user := h.svc.Auth.GetCurrentUser(ctx)

	log.Printf("current user: %v", user.ID)

	return rest.SuccessResponse(ctx, "category endpoint", nil)
}

func (h CatalogHandler) EditCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "edit category endpoint", nil)
}

func (h CatalogHandler) DeleteCategory(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "delete category endpoint", nil)
}

// Products

func (h CatalogHandler) CreateProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Create product endpoint", nil)
}

func (h CatalogHandler) GetProducts(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Get productS endpoint", nil)
}

func (h CatalogHandler) GetProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Get product endpoint", nil)
}

func (h CatalogHandler) EditProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Edit product endpoint", nil)
}

func (h CatalogHandler) UpdateProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Update product endpoint", nil)
}

func (h CatalogHandler) DeleteProduct(ctx *fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Delete product endpoint", nil)
}
