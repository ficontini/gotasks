package db

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const projectColl = "projects"

type ProjectStore interface {
	GetProjectByID(context.Context, string) (*types.Project, error)
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
func (s *MongoProjectStore) GetProjectByID(ctx context.Context, id string) (*types.Project, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var project *types.Project
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&project); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return project, nil
}
