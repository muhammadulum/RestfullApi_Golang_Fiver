package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"learn_restful_api_golang/internal/domain"
	"learn_restful_api_golang/internal/model"
)

type UserHandler struct {
	usecase domain.AuthUseCase
}

func NewUserHandler(u domain.AuthUseCase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) RegisterRoutes(r fiber.Router) {
	r.Post("/login", h.Login)
	r.Post("/register", h.Register)
	r.Post("/refresh", h.Refresh)
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body model.LoginRequest true "Login data"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /login [post]



func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	token, refresh, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"access_token": token, "refresh_token": refresh})
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req model.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	err := h.usecase.Register(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "user registered"})
}

func (h *UserHandler) Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}
	token, err := h.usecase.RefreshToken(body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}
	return c.JSON(fiber.Map{"access_token": token})
}

// ...existing code...

func (h *UserHandler) ProtectedExample(c *fiber.Ctx) error {
    claims := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)
    role := claims["role"].(string)
    if role != "admin" {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden access"})
    }
    return c.JSON(fiber.Map{"message": "you are authenticated as admin"})
}