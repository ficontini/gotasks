package main

import (
	"context"
	"flag"
	"log"

	"github.com/ficontini/gotasks/api"
	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	config = fiber.Config{
		ErrorHandler: api.ErrorHandler,
	}
)

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore = db.NewMongoUserStore(client)
		taskStore = db.NewMongoTaskStore(client)
		store     = db.Store{
			User:    userStore,
			Task:    taskStore,
			Project: db.NewMongoProjectStore(client, taskStore),
		}
		taskService    = service.NewTaskService(taskStore)
		projectService = service.NewProjectService(store)
		userService    = service.NewUserService(userStore)
		taskHandler    = api.NewTaskHandler(taskService)
		userHandler    = api.NewUserHandler(*userService)
		authHandler    = api.NewAuthHandler(userStore)
		projectHandler = api.NewProjectHandler(projectService)
		app            = fiber.New(config)
		auth           = app.Group("/api")
		apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
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

	apiv1.Post("project/", projectHandler.HandlePostProject)
	apiv1.Post("project/:id/task", projectHandler.HandleAddTaskToProject)
	apiv1.Get("project/:id/task", projectHandler.HandleGetTasks)

	log.Fatal(app.Listen(*listenAddr))

}
