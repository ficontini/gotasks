package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/ficontini/gotasks/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoTaskStore(client *mongo.Client) *MongoTaskStore {
	return &MongoTaskStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(taskColl),
	}
}
func (s *MongoTaskStore) InsertTask(ctx context.Context, task *data.Task) (*data.Task, error) {
	res, err := s.coll.InsertOne(ctx, task)
	if err != nil {
		return nil, err
	}
	task.ID = data.CreateIDFromObjectID(res.InsertedID.(primitive.ObjectID))
	return task, nil
}

func (s *MongoTaskStore) Update(ctx context.Context, id data.ID, params data.UpdateTaskParams) error {
	update := bson.M{"$set": params.ToMap()}
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

func (s *MongoTaskStore) UpdateTaskProjects(ctx context.Context, filter Map, update Map) error {
	res, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrorNotFound
	}
	return nil
}

func (s *MongoTaskStore) Delete(ctx context.Context, id data.ID) error {
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
func (s *MongoTaskStore) GetTaskByID(ctx context.Context, id data.ID) (*data.Task, error) {
	oid, err := id.ObjectID()
	if err != nil {
		return nil, err
	}
	var task *data.Task
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&task); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return task, nil
}
func (s *MongoTaskStore) GetTasks(ctx context.Context, filter Map, pagination *Pagination) ([]*data.Task, error) {
	return s.getTasks(ctx, filter, pagination)
}
func (s *MongoTaskStore) GetTasksByProject(ctx context.Context, id data.ID, pagination *Pagination) ([]*data.Task, error) {
	oid, err := id.ObjectID()
	if err != nil {
		return nil, err
	}
	return s.getTasks(ctx, Map{"projects": oid}, pagination)
}
func (s *MongoTaskStore) getTasks(ctx context.Context, filter Map, pagination *Pagination) ([]*data.Task, error) {
	opts := pagination.getOptions()
	fmt.Println("---", pagination)
	cur, err := s.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var tasks []*data.Task
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, err
}
