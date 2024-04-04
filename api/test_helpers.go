package api

import (
	"context"
	"log"
	"testing"

	"github.com/ficontini/gotasks/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

const testdburi = "mongodb://localhost:27017"

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	taskStore := db.NewMongoTaskStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			Task:    taskStore,
			User:    db.NewMongoUserStore(client),
			Project: db.NewMongoProjectStore(client, taskStore),
		},
	}
}

func checkStatusCode(t *testing.T, expected, actual int) {
	if actual != expected {
		t.Fatalf("expected %d status code, but got %d", expected, actual)
	}
}
