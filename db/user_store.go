package db

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/data"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserByID(context.Context, data.ID) (*data.User, error)
	GetUserByEmail(context.Context, string) (*data.User, error)
	InsertUser(context.Context, *data.User) (*data.User, error)
	EnableUser(context.Context, data.ID) error
}

type MongoUserStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*data.User, error) {
	var user *data.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *MongoUserStore) GetUserByID(ctx context.Context, id data.ID) (*data.User, error) {
	oid, err := id.ObjectID()
	if err != nil {
		return nil, err
	}
	var user *data.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *MongoUserStore) InsertUser(ctx context.Context, user *data.User) (*data.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = data.CreateIDFromObjectID(res.InsertedID.(primitive.ObjectID))
	return user, nil
}
func (s *MongoUserStore) EnableUser(ctx context.Context, id data.ID) error {
	oid, err := id.ObjectID()
	if err != nil {
		return err
	}
	update := bson.M{"$set": bson.M{"enabled": true}}
	res, err := s.coll.UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrorNotFound
	}
	return nil
}
