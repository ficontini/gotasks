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
	store, err := db.NewDynamoDBStore()
	if err != nil {
		panic(err)
	}
	var (
		handler = api.NewHandler(service.NewService(store))
		svc     = service.NewProjectService(store)
		app     = fiber.New(config)
		auth    = app.Group("/api")
		apiv1   = app.Group("/api/v1", api.JWTAuthentication(service.NewAuthService(store)))
		admin   = apiv1.Group("/admin", api.AdminAuth)
	)
	svc = service.NewLogMiddleware(svc)
	projectHandler := api.NewProjectHandler(svc)
	auth.Post("/auth", handler.Auth.HandleAuthenticate)

	auth.Post("/user", handler.User.HandlePostUser)
	apiv1.Post("/user/reset-password", handler.User.HandleResetPassword)
	apiv1.Get("/user", handler.User.HandleGetUser)

	apiv1.Get("/task/all", handler.Task.HandleGetTasks)
	apiv1.Get("/task", handler.Task.HandleGetUserTasks)
	apiv1.Post("/task", handler.Task.HandlePostTask)
	apiv1.Get("/task/:id", handler.Task.HandleGetTask)
	apiv1.Post("/task/:id/assign", handler.Task.HandleAssignTaskToSelf)
	apiv1.Post("/task/:id/complete", handler.Task.HandleCompleteTask)
	apiv1.Put("/task/:id/due-date", handler.Task.HandlePutDueDateTask)
	admin.Post("/task/:id/assign", handler.Task.HandleAssignTaskToUser)
	admin.Delete("/task/:id", handler.Task.HandleDeleteTask)
	admin.Get("/task", handler.Task.HandleGetTasks)

	admin.Get("/user", handler.User.HandleGetUsers)
	admin.Get("/user/:id", handler.User.HandleAdminGetUser)
	admin.Put("/user/:id/enable", handler.User.HandleEnableUser)
	admin.Put("/user/:id/disable", handler.User.HandleDisableUser)

	apiv1.Post("/project", projectHandler.HandlePostProject)
	apiv1.Get("/project/:id", projectHandler.HandleGetProject)
	apiv1.Post("/project/:id/task", projectHandler.HandlePostTask)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listenAddr))
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
