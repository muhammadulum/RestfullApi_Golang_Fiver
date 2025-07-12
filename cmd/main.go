package main

import (
	"learn_restful_api_golang/config"
	"learn_restful_api_golang/internal/handler"
	"learn_restful_api_golang/internal/repository"
	"learn_restful_api_golang/internal/usecase"
	"learn_restful_api_golang/internal/domain"
	"github.com/gofiber/fiber/v2"
	"learn_restful_api_golang/pkg/middleware"

//   _ "learn_restful_api_golang/docs" // import for swagger docs

//     "github.com/gofiber/fiber/v2"
//     fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	db := config.InitPostgres()
	defer config.CloseDB(db)

	// Auto migrate model
	db.AutoMigrate(&domain.User{})

	r := fiber.New()

	repo := repository.NewUserRepository(db)
	usecase := usecase.NewUserUseCase(repo)
	handler := handler.NewUserHandler(usecase)

	api := r.Group("/api")
	handler.RegisterRoutes(api)

	// Protected route example
	api.Get("/profile", middleware.JWTProtected("user"), handler.ProtectedExample)

	// Swagger endpoint


	r.Listen(":3000")
}