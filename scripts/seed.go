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
	store := seedDynamoDB()
	for i := 0; i < 10; i++ {
		title := fmt.Sprintf("task%d", i)
		description := fmt.Sprintf("description of task%d", i)
		fixtures.AddTask(store, title, description, time.Now().AddDate(0, 0, rand.Intn(10)), rand.Intn(2) == 0)
	}
	fixtures.AddUser(store, "james", "foobaz", "supersecurepassword", false, true)
	fixtures.AddUser(store, "admin", "foobaz", "supersecurepassword", true, true)
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
