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
	task.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return task, nil
}

func (s *MongoTaskStore) SetTaskAsComplete(ctx context.Context, id string, params types.TaskCompletionRequest) error {
	update := bson.M{"$set": params.ToMap()}
	return s.update(ctx, id, update)
}
func (s *MongoTaskStore) SetTaskAssignee(ctx context.Context, id string, req types.TaskAssignmentRequest) error {
	oid, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return err
	}
	params := req.ToMap()
	params["assignedTo"] = oid
	update := bson.M{"$set": params}
	return s.update(ctx, id, update)
}
func (s *MongoTaskStore) update(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
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

func (s *MongoTaskStore) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
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
func (s *MongoTaskStore) GetTaskByID(ctx context.Context, id string) (*types.Task, error) {
	oid, err := primitive.ObjectIDFromHex(id)
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
func (s *MongoTaskStore) GetTasks(ctx context.Context, filter types.Filter, pagination *Pagination) ([]*types.Task, error) {
	opts := pagination.getOptions()
	cur, err := s.coll.Find(ctx, filter.Apply(), opts)
	if err != nil {
		return nil, err
	}
	var tasks []*types.Task
	if err := cur.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, err
}
