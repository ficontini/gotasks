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
	TransactAddTask(context.Context, []*UpdateAction) error
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
func (s *MongoProjectStore) InsertProject(ctx context.Context, project *types.Project) (*types.Project, error) {
	res, err := s.coll.InsertOne(ctx, project)
	if err != nil {
		return nil, ErrInvalidID
	}
	project.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return project, nil
}
func (s *MongoProjectStore) GetProjectByID(ctx context.Context, id string) (*types.Project, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidID
	}
	var project *types.Project
	if err := s.coll.FindOne(ctx, bson.M{mongoIDField: oid}).Decode(&project); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return project, nil
}
func (s *MongoProjectStore) Update(ctx context.Context, id string, params Update) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update, err := params.ToBSON()
	if err != nil {
		return err
	}
	res, err := s.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrorNotFound
	}
	return nil
}

// TODO: Review
func (s *MongoProjectStore) TransactAddTask(ctx context.Context, actions []*UpdateAction) error {
	session, err := s.client.StartSession()
	defer session.EndSession(ctx)
	if err != nil {
		return err
	}
	session.StartTransaction()
	for _, action := range actions {
		switch action.TableName {
		case taskColl:
			err = s.TaskStore.Update(ctx, action.ID, action.Params)
			if err != nil {
				session.AbortTransaction(ctx)
				return err
			}
		default:
			err = s.Update(ctx, action.ID, action.Params)
			if err != nil {
				session.AbortTransaction(ctx)
				return err
			}
		}
	}
	err = session.CommitTransaction(ctx)
	return err
}
