package main

import (
	"log"
	"os"

	"github.com/ficontini/gotasks/api"
	"github.com/ficontini/gotasks/db"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	config = fiber.Config{
		ErrorHandler: api.ErrorHandler,
	}
	loggerConfig = logger.Config{
		Next: func(c *fiber.Ctx) bool {
			logrus.WithFields(logrus.Fields{
				"method": c.Method(),
				"path":   c.Path(),
			}).Info("Request")
			return true
		},
	}
)

func main() {
	var (
		app = fiber.New(config)
	)
	cfg, err := db.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	app.Use(logger.New(loggerConfig))
	MakeRoutes(app, cfg.Store)
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listenAddr))
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
