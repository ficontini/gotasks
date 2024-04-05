package db

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const projectColl = "projects"

type ProjectStore interface {
	GetProjectByID(context.Context, string) (*data.Project, error)
	InsertProject(context.Context, *data.Project) (*data.Project, error)
}

type MongoProjectStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	TaskStore
}

func NewMongoProjectStore(client *mongo.Client, taskStore TaskStore) *MongoProjectStore {
	return &MongoProjectStore{
		client:    client,
		coll:      client.Database(DBNAME).Collection(projectColl),
		TaskStore: taskStore,
	}
}
func (s *MongoProjectStore) InsertProject(ctx context.Context, project *data.Project) (*data.Project, error) {
	res, err := s.coll.InsertOne(ctx, project)
	if err != nil {
		return nil, err
	}
	project.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return project, nil
}
func (s *MongoProjectStore) GetProjectByID(ctx context.Context, id string) (*data.Project, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var project *data.Project
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&project); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return project, nil
}
