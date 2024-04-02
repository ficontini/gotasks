package db

import (
	"context"
	"errors"

	"github.com/ficontini/gotasks/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userColl = "users"

type UserStore interface {
	GetUserByID(context.Context, types.ID) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	InsertUser(context.Context, *types.User) (*types.User, error)
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

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	var user *types.User
	if err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *MongoUserStore) GetUserByID(ctx context.Context, id types.ID) (*types.User, error) {
	oid, err := id.ObjectID()
	if err != nil {
		return nil, err
	}
	var user *types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrorNotFound
		}
		return nil, err
	}
	return user, nil
}
func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = types.CreateIDFromObjectID(res.InsertedID.(primitive.ObjectID))
	return user, nil
}
