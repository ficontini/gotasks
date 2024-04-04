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
	GetProjectByID(context.Context, data.ID) (*data.Project, error)
	InsertProject(context.Context, *data.Project) (*data.Project, error)
	UpdateProjectTasks(context.Context, data.ID, data.ID) error
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
	project.ID = data.CreateIDFromObjectID(res.InsertedID.(primitive.ObjectID))
	return project, nil
}
func (s *MongoProjectStore) GetProjectByID(ctx context.Context, id data.ID) (*data.Project, error) {
	oid, err := id.ObjectID()
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

// TODO: review
func (s *MongoProjectStore) UpdateProjectTasks(ctx context.Context, projectID data.ID, taskID data.ID) error {
	oprojectID, err := projectID.ObjectID()
	if err != nil {
		return err
	}
	otaskID, err := taskID.ObjectID()
	if err != nil {
		return err
	}
	update := Map{"$push": Map{"tasks": otaskID}}
	res, err := s.coll.UpdateByID(ctx, oprojectID, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrorNotFound
	}
	updateTask := Map{"$push": Map{"projects": oprojectID}}
	err = s.TaskStore.UpdateTaskProjects(ctx, Map{"_id": otaskID}, updateTask)
	return err
}
