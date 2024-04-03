package db

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

// TODO: review
func NewMongoTaskStore(client *mongo.Client) *MongoTaskStore {
	return &MongoTaskStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(taskColl),
	}
}
func (s *MongoTaskStore) InsertTask(ctx context.Context, task *types.Task) (*types.Task, error) {
	res, err := s.coll.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	task.ID = types.CreateIDFromObjectID(res.InsertedID.(primitive.ObjectID))
	return task, nil
}
func (s *MongoTaskStore) Update(ctx context.Context, id types.ID, params types.UpdateTaskParams) error {
	update := bson.M{"$set": params.ToMap()}
	return s.update(ctx, id, update)
}
func (s *MongoTaskStore) UpdateTaskProjects(ctx context.Context, id types.ID, projectID types.ID) error {
	oprojectID, err := projectID.ObjectID()
	if err != nil {
		return err
	}
	update := bson.M{"$push": bson.M{"projects": oprojectID}}
	return s.update(ctx, id, update)
}
func (s *MongoTaskStore) update(ctx context.Context, id types.ID, update bson.M) error {
	oid, err := id.ObjectID()
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
func (s *MongoTaskStore) Delete(ctx context.Context, id types.ID) error {
	oid, err := id.ObjectID()
	if err != nil {
		return err
	}
	res, err := s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrorNotFound
	}
	return nil
}
func (s *MongoTaskStore) GetTaskByID(ctx context.Context, id types.ID) (*types.Task, error) {
	oid, err := id.ObjectID()
	if err != nil {
		return nil, err
	}
	var task *types.Task
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&task); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return task, nil
}
func (s *MongoTaskStore) GetTasks(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Task, error) {
	return s.getTasks(ctx, filter, pagination)
}
func (s *MongoTaskStore) GetTasksByProject(ctx context.Context, id types.ID, pagination *Pagination) ([]*types.Task, error) {
	oid, err := id.ObjectID()
	if err != nil {
		return nil, err
	}
	return s.getTasks(ctx, Map{"projects": oid}, pagination)
}
func (s *MongoTaskStore) getTasks(ctx context.Context, filter Map, pagination *Pagination) ([]*types.Task, error) {
	opts := pagination.getOptions()
	cur, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var tasks []*types.Task
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, err
}
