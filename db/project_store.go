package db

import (
	"context"

	"github.com/ficontini/gotasks/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const projectColl = "projects"

type ProjectStore interface {
	InsertProject(context.Context, *types.Project) (*types.Project, error)
}

type MongoProjectStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoProjectStore(client *mongo.Client) *MongoProjectStore {
	return &MongoProjectStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(projectColl),
	}
}
func (s *MongoProjectStore) InsertProject(ctx context.Context, project *types.Project) (*types.Project, error) {
	res, err := s.coll.InsertOne(ctx, project)
	if err != nil {
		return nil, err
	}
	project.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return project, nil
}
