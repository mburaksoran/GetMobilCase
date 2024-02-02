package client_api

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/mburaksoran/GetMobilCase/client_api/internal/app/config"
	"golang.org/x/oauth2"
)

type googleAuthHandler struct {
	cfg          *config.AppConfig
	googleConfig oauth2.Config
}
type GoogleAuthInterface interface {
	GoogleLogin(c *fiber.Ctx) error
	GoogleCallback(c *fiber.Ctx) error
}

func NewGoogleAuthHandler(cfg *config.AppConfig) GoogleAuthInterface {
	gcfg := config.GoogleConfig(cfg)
	return &googleAuthHandler{cfg: cfg, googleConfig: gcfg}
}

func (g *googleAuthHandler) GoogleLogin(c *fiber.Ctx) error {

	url := g.googleConfig.AuthCodeURL("randomstate")

	c.Status(fiber.StatusSeeOther)
	c.Redirect(url)
	return c.JSON(url)
}

func (g *googleAuthHandler) GoogleCallback(c *fiber.Ctx) error {

	state := c.Query("state")
	if state != "randomstate" {
		return c.SendString("States don't Match!!")
	}

	code := c.Query("code")

	token, err := g.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return c.SendString("Code-Token Exchange Failed")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.SendString("User Data Fetch Failed")
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.SendString("JSON Parsing Failed")
	}

	return c.SendString(string(token.AccessToken))
}
