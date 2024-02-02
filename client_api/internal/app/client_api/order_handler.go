package client_api

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/models"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/domain/service"
)

type orderHandler struct {
	OrderService service.OrderServiceInterface
	googleConfig oauth2.Config
}

type OrderHandlerInterface interface {
	CreateOrder(c *fiber.Ctx) error
}

func NewOrderHandler(orderService service.OrderServiceInterface) OrderHandlerInterface {
	return &orderHandler{OrderService: orderService}
}
func (o *orderHandler) CreateOrder(c *fiber.Ctx) error {
	//isVerified, err := o.verifyToken(c)
	//if err != nil {
	//	return c.SendStatus(fiber.StatusInternalServerError)
	//}
	//if !isVerified {
	//	return c.SendStatus(fiber.StatusUnauthorized)
	//}
	var order *models.Order
	err := c.BodyParser(&order)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = o.OrderService.CreateOrder(c.Context(), order)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (g *orderHandler) verifyToken(c *fiber.Ctx) (bool, error) {
	code := c.Query("code")
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + code)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var userInfo map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userInfo)
	if err != nil {
		log.Printf("Failed to decode userinfo: %s", err.Error())
		return false, err
	}
	return true, nil
}
