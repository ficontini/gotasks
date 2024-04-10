package db

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Update interface {
	ToBSON() bson.M
}

type StatusUpdater struct {
	Enabled bool
}

func (u StatusUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"enabled": u.Enabled},
	}
}

type PasswordUpdater struct {
	EncryptedPassword string
}

func (u PasswordUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"encryptedPassword": u.EncryptedPassword},
	}
}

type TaskCompleteUpdater struct {
	Completed bool
}

func (u TaskCompleteUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"completed": u.Completed},
	}
}

type TaskAssignationUpdater struct {
	AssignedTo interface{}
}

func NewTaskAssignationUpdater(assignedTo string) (*TaskAssignationUpdater, error) {
	oid, err := primitive.ObjectIDFromHex(assignedTo)
	if err != nil {
		return nil, err
	}
	return &TaskAssignationUpdater{
		AssignedTo: oid,
	}, nil
}
func (u TaskAssignationUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"assignedTo": u.AssignedTo},
	}
}
