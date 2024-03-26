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

	taskHandler := api.NewTaskHandler(db.NewMongoTaskStore(client))

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/task", taskHandler.HandleGetTasks)
	apiv1.Post("/task", taskHandler.HandlePostTask)
	apiv1.Get("/task/:id", taskHandler.HandleGetTask)
	apiv1.Post("/task/:id/complete", taskHandler.HandleCompleteTask)
	apiv1.Delete("/task/:id", taskHandler.HandleDeleteTask)

	log.Fatal(app.Listen(*listenAddr))
}
