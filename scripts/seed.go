package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ficontini/gotasks/db"
	"github.com/ficontini/gotasks/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.Background()

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		Task: db.NewMongoTaskStore(client),
		User: db.NewMongoUserStore(client),
	}
	for i := 0; i < 100; i++ {
		title := fmt.Sprintf("task%d", i)
		description := fmt.Sprintf("description of task%d", i)
		fixtures.AddTask(store, title, description, time.Now().AddDate(0, 0, rand.Intn(10)), rand.Intn(2) == 0)
	}
	fixtures.AddUser(store, "james", "foobaz", "supersecurepassword", false, true)
	fixtures.AddUser(store, "admin", "foobaz", "supersecurepassword", true, true)
	fmt.Println("seeding the database")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	if err := db.SetupMongoConfigFromEnv(); err != nil {
		log.Fatal(err)
	}

}
