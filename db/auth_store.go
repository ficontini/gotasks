package db

import (
	"context"

	"github.com/ficontini/gotasks/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const authColl = "auths"

type AuthStore interface {
	Insert(context.Context, *types.Auth) (*types.Auth, error)
	Get(context.Context, *types.AuthFilter) (*types.Auth, error)
	Delete(context.Context, *types.AuthFilter) error
}

type MongoAuthStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoAuthStore(client *mongo.Client) *MongoAuthStore {
	return &MongoAuthStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(authColl),
	}
}

func (s *MongoAuthStore) Insert(ctx context.Context, auth *types.Auth) (*types.Auth, error) {
	mAuth, err := newMongoAuth(auth)
	if err != nil {
		return nil, err
	}
	res, err := s.coll.InsertOne(ctx, mAuth)
	if err != nil {
		return nil, err
	}
	auth.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return auth, err
}

func (s *MongoAuthStore) Get(ctx context.Context, params *types.AuthFilter) (*types.Auth, error) {
	filter, err := NewMongoAuthFilter(params)
	if err != nil {
		return nil, err
	}
	var auth *types.Auth
	if err := s.coll.FindOne(ctx, filter).Decode(&auth); err != nil {
		return nil, err
	}
	return auth, err
}
func (s *MongoAuthStore) Delete(ctx context.Context, params *types.AuthFilter) error {
	filter, err := NewMongoAuthFilter(params)
	if err != nil {
		return err
	}
	res, err := s.coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrorNotFound
	}
	return nil
}

type MongoAuth struct {
	ID       string             `bson:"_id,omitempty"`
	UserID   primitive.ObjectID `bson:"userID"`
	AuthUUID string             `bson:"authUUID"`
}

func newMongoAuth(auth *types.Auth) (*MongoAuth, error) {
	oid, err := primitive.ObjectIDFromHex(auth.UserID)
	if err != nil {
		return nil, err
	}
	return &MongoAuth{
		UserID:   oid,
		AuthUUID: auth.AuthUUID,
	}, nil
}

type MongoAuthFilter struct {
	UserID   primitive.ObjectID `bson:"userID"`
	AuthUUID string             `bson:"authUUID"`
}

func NewMongoAuthFilter(params *types.AuthFilter) (*MongoAuthFilter, error) {
	oid, err := primitive.ObjectIDFromHex(params.UserID)
	if err != nil {
		return nil, err
	}
	return &MongoAuthFilter{
		UserID:   oid,
		AuthUUID: params.AuthUUID,
	}, nil
}
