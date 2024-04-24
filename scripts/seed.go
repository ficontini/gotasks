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
)

func main() {
	store := seedMongo()
	for i := 0; i < 50; i++ {
		title := fmt.Sprintf("task%d", i)
		description := fmt.Sprintf("description of task%d", i)
		fixtures.AddTask(store, title, description, time.Now().AddDate(0, 0, rand.Intn(10)), rand.Intn(2) == 0)
	}
	fixtures.AddUser(store, "james", "foobaz", "supersecurepassword", false, true)
	fixtures.AddUser(store, "admin", "foobaz", "supersecurepassword", true, true)
	// fixtures.AddUser(store, "luca", "foobaz", "supersecurepassword", false, true)
	// fixtures.AddUser(store, "frank", "foobaz", "supersecurepassword", false, true)
	// fixtures.AddUser(store, "doni", "foobaz", "supersecurepassword", false, true)
	// fixtures.AddUser(store, "toni", "foobaz", "supersecurepassword", false, true)
	// fixtures.AddUser(store, "antonio", "foobaz", "supersecurepassword", false, true)
	// fixtures.AddUser(store, "ferland", "foobaz", "supersecurepassword", false, true)
	// fixtures.AddUser(store, "admin2", "foobaz", "supersecurepassword", true, true)
	// fixtures.AddUser(store, "enrique", "foobaz", "supersecurepassword", false, true)
	fmt.Println("seeding the database")
}
func seedDynamoDB() *db.Store {
	client, err := db.NewDynamoDBClient()
	if err != nil {
		log.Fatal(err)
	}
	return &db.Store{
		Task: db.NewDynamoDBTaskStore(client),
		User: db.NewDynamoDBUserStore(client),
	}
}
func seedMongo() *db.Store {
	client, err := db.NewMongoClient()
	if err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		Task: db.NewMongoTaskStore(client),
		User: db.NewMongoUserStore(client),
	}
	if err := store.Task.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}

	if err := store.User.Drop(context.Background()); err != nil {
		log.Fatal(err)
	}
	return store
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
