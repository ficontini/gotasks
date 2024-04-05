package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ficontini/gotasks/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MongoTestEndpoint = "MONGO_DB_TEST_URI"

var (
	testdburi string
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func setup(t *testing.T) *testdb {
	if err := Init(); err != nil {
		log.Fatal(err)
	}
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

func Init() error {
	testdburi = os.Getenv(MongoTestEndpoint)
	if testdburi == "" {
		return fmt.Errorf("%s env variable not set", MongoTestEndpoint)
	}
	db.DBNAME = os.Getenv(db.MongoDBNameEnvName)
	if db.DBNAME == "" {
		return fmt.Errorf("%s env variable not set", db.MongoDBNameEnvName)
	}
	return nil
}

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
}
