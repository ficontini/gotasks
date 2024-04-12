package db

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Update interface {
	ToBSON() bson.M
	ToExpression() expression.UpdateBuilder
}

type StatusUpdater struct {
	Enabled bool
}

func (u StatusUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"enabled": u.Enabled},
	}
}
func (u StatusUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name("enabled"), expression.Value(u.Enabled))
}

type PasswordUpdater struct {
	EncryptedPassword string
}

func (u PasswordUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"encryptedPassword": u.EncryptedPassword},
	}
}
func (u PasswordUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name("encryptedPassword"), expression.Value(u.EncryptedPassword))
}

type TaskCompleteUpdater struct {
	Completed bool
}

func (u TaskCompleteUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"completed": u.Completed},
	}
}
func (u TaskCompleteUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name("completed"), expression.Value(u.Completed))
}

type TaskAssignationUpdater struct {
	AssignedTo string
}

func (u TaskAssignationUpdater) ToBSON() bson.M {
	oid, err := primitive.ObjectIDFromHex(u.AssignedTo)
	if err != nil {
		log.Fatal(err)
	}
	return bson.M{
		"$set": bson.M{"assignedTo": oid},
	}
}
func (u TaskAssignationUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name("assignedTo"), expression.Value(u.AssignedTo))
}
