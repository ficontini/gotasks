package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/ficontini/gotasks/db"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	MongoTestEndpoint = "MONGO_DB_TEST_URI"
	EnvFile           = "../.env"
)

var (
	testdburi string
)

func setup(t *testing.T) TestDB {
	return setupTestDynamoDB(t)
}

type TestDB interface {
	Store() *db.Store
	teardown(*testing.T)
}

type TestMongoDB struct {
	client *mongo.Client
	store  *db.Store
}

func setupTestMongoDB(t *testing.T) *TestMongoDB {
	if err := InitMongoTestDB(); err != nil {
		log.Fatal(err)
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(testdburi))
	if err != nil {
		log.Fatal(err)
	}
	taskStore := db.NewMongoTaskStore(client)
	return &TestMongoDB{
		client: client,
		store: &db.Store{
			Task:    taskStore,
			User:    db.NewMongoUserStore(client),
			Project: db.NewMongoProjectStore(client, taskStore),
			Auth:    db.NewMongoAuthStore(client),
		},
	}
}
func (tdb *TestMongoDB) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
func (tdb *TestMongoDB) Store() *db.Store {
	return tdb.store
}

func InitMongoTestDB() error {
	if err := godotenv.Load(EnvFile); err != nil {
		log.Fatal(err)
		return err
	}
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

type TestDynamoDB struct {
	client *dynamodb.Client
	store  *db.Store
}

// TODO:
func (tdb *TestDynamoDB) teardown(t *testing.T) {

}
func setupTestDynamoDB(t *testing.T) *TestDynamoDB {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
	if err := db.SetupDynamoDBConfigFromEnv(); err != nil {
		log.Fatal(err)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithSharedConfigProfile(db.PROFILE))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	client := dynamodb.NewFromConfig(cfg)
	taskStore := db.NewDynamoDBTaskStore(client)
	return &TestDynamoDB{
		client: client,
		store: &db.Store{
			Auth:    db.NewDynamoDBAuthStore(client),
			User:    db.NewDynamoDBUserStore(client),
			Task:    taskStore,
			Project: db.NewDynamoDBProjectStore(client, taskStore),
		},
	}
}
func (tdb *TestDynamoDB) Store() *db.Store {
	return tdb.store
}
