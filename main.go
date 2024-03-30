package main

import (
	"context"
	"flag"
	"log"

	"github.com/ficontini/gotasks/api"
	"github.com/ficontini/gotasks/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":3000", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	var (
		userStore    = db.NewMongoUserStore(client)
		taskStore    = db.NewMongoTaskStore(client)
		projectStore = db.NewMongoProjectStore(client)
		store        = db.Store{
			User:    userStore,
			Task:    taskStore,
			Project: projectStore,
		}
		taskHandler    = api.NewTaskHandler(taskStore)
		userHandler    = api.NewUserHandler(userStore)
		authHandler    = api.NewAuthHandler(userStore)
		projectHandler = api.NewProjectHandler(&store)
		app            = fiber.New(config)
		apiv1          = app.Group("/api/v1")
		apiv1Task      = apiv1.Group("/task", api.JWTAuthentication(userStore))
		apiv1Project   = apiv1.Group("/project", api.JWTAuthentication(userStore))
	)

	apiv1.Post("/auth", authHandler.HandleAuthenticate)
	apiv1.Post("/user", userHandler.HandlePostUser)

	//TODO: review groups
	apiv1Task.Get("/", taskHandler.HandleGetTasks)
	apiv1Task.Post("/", taskHandler.HandlePostTask)
	apiv1Task.Get("/:id", taskHandler.HandleGetTask)
	apiv1Task.Post("/:id/complete", taskHandler.HandleCompleteTask)
	apiv1Task.Delete("/:id", taskHandler.HandleDeleteTask)

	apiv1Project.Post("/", projectHandler.HandlePostProject)

	log.Fatal(app.Listen(*listenAddr))

}
