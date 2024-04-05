package main

import (
	"log"
	"os"

	"github.com/ficontini/gotasks/api"
	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	config = fiber.Config{
		ErrorHandler: api.ErrorHandler,
	}
)

func main() {
	store, err := db.NewMongoStore()
	if err != nil {
		panic(err)
	}

	var (
		taskHandler    = api.NewTaskHandler(service.NewTaskService(store.Task))
		userHandler    = api.NewUserHandler(*service.NewUserService(store.User))
		authHandler    = api.NewAuthHandler(store.User)
		projectHandler = api.NewProjectHandler(service.NewProjectService(*store))
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(store.User))
		admin          = apiv1.Group("/admin", api.AdminAuth)
	)

	auth.Post("/auth", authHandler.HandleAuthenticate)
	auth.Post("/user", userHandler.HandlePostUser)

	//admin tasks
	admin.Post("/user/:id/enable", userHandler.HandleEnableUser)
	admin.Delete("task/:id", taskHandler.HandleDeleteTask)

	apiv1.Get("/task", taskHandler.HandleGetTasks)
	apiv1.Post("/task", taskHandler.HandlePostTask)
	apiv1.Get("task/:id", taskHandler.HandleGetTask)
	apiv1.Post("task/:id/complete", taskHandler.HandleCompleteTask)

	apiv1.Post("/project", projectHandler.HandlePostProject)
	apiv1.Get("/project/:id", projectHandler.HandleGetProject)
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listenAddr))
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
