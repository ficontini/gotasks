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

type UserGetter interface {
	GetUsers(context.Context, Filter, Pagination) ([]*types.User, error)
	GetUserByID(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
}
type UserStore interface {
	UserGetter
	InsertUser(context.Context, *types.User) (*types.User, error)
	Update(context.Context, string, Update) error
	Dropper
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
func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
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
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()
	return user, nil
}
func (s *MongoUserStore) Update(ctx context.Context, id string, update Update) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	res, err := s.coll.UpdateByID(ctx, oid, update.ToBSON())
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return ErrorNotFound
	}
	return nil
}
func (s *MongoUserStore) GetUsers(ctx context.Context, filter Filter, pagination Pagination) ([]*types.User, error) {
	opts := pagination.getOptions()
	cur, err := s.coll.Find(ctx, filter.ToBSON(), opts)
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, err
}
func (s *MongoUserStore) Drop(ctx context.Context) error {
	return s.coll.Drop(ctx)
}
