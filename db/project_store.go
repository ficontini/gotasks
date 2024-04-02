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
	GetProjectByID(context.Context, types.ID) (*types.Project, error)
	InsertProject(context.Context, *types.Project) (*types.Project, error)
	Updater
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
	project.ID = types.CreateIDFromObjectID(res.InsertedID.(primitive.ObjectID))
	return project, nil
}
func (s *MongoProjectStore) GetProjectByID(ctx context.Context, id types.ID) (*types.Project, error) {
	oid, err := id.ObjectID()
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
func (s *MongoProjectStore) Update(ctx context.Context, id types.ID, update Map) error {
	oid, err := id.ObjectID()
	if err != nil {
		return err
	}
	if pushUpdate, ok := update["$push"].(Map); ok {
		taskID, err := pushUpdate["tasks"].(types.ID).ObjectID()
		if err != nil {
			return err
		}
		pushUpdate["tasks"] = taskID
	}
	_, err = s.coll.UpdateByID(ctx, oid, update)
	return err
}
