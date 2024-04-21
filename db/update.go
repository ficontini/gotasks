package db

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Update interface {
	ToBSON() (bson.M, error)
	ToExpression() expression.UpdateBuilder
}

type StatusUpdater struct {
	Enabled bool
}

func (u StatusUpdater) ToBSON() (bson.M, error) {
	return bson.M{
		"$set": bson.M{enabledField: u.Enabled},
	}, nil
}
func (u StatusUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(enabledField), expression.Value(u.Enabled))
}

type PasswordUpdater struct {
	EncryptedPassword string
}

func (u PasswordUpdater) ToBSON() (bson.M, error) {
	return bson.M{
		"$set": bson.M{encryptedPasswordField: u.EncryptedPassword},
	}, nil
}
func (u PasswordUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(encryptedPasswordField), expression.Value(u.EncryptedPassword))
}

type TaskDueDateUpdater struct {
	DueDate time.Time
}

func (u TaskDueDateUpdater) ToBSON() (bson.M, error) {
	return bson.M{
		"$set": bson.M{dueDateField: u.DueDate},
	}, nil
}
func (u TaskDueDateUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(dueDateField), expression.Value(u.DueDate))
}

type TaskCompleteUpdater struct {
	Completed bool
}

func (u TaskCompleteUpdater) ToBSON() (bson.M, error) {
	return bson.M{
		"$set": bson.M{completedField: u.Completed},
	}, nil
}
func (u TaskCompleteUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(completedField), expression.Value(u.Completed))
}

type TaskAssignationUpdater struct {
	AssignedTo string
}

func (u TaskAssignationUpdater) ToBSON() (bson.M, error) {
	oid, err := primitive.ObjectIDFromHex(u.AssignedTo)
	if err != nil {
		return nil, err
	}
	return bson.M{
		"$set": bson.M{assignedToField: oid},
	}, nil
}
func (u TaskAssignationUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(assignedToField), expression.Value(u.AssignedTo))
}

type AddTaskToProjectUpdater struct {
	TaskID string
}

func (u AddTaskToProjectUpdater) ToBSON() (bson.M, error) {
	oid, err := primitive.ObjectIDFromHex(u.TaskID)
	if err != nil {
		return nil, err
	}
	return bson.M{
		"$push": bson.M{tasksField: oid},
	}, nil
}
func (u AddTaskToProjectUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(tasksField), expression.ListAppend(expression.Name(tasksField), expression.Value([]string{u.TaskID})))
}

type TaskProjectIDUpdater struct {
	ProjectID string
}

func (u TaskProjectIDUpdater) ToBSON() (bson.M, error) {
	oid, err := primitive.ObjectIDFromHex(u.ProjectID)
	if err != nil {
		return nil, err
	}
	return bson.M{"$set": bson.M{projectIDField: oid}}, nil
}
func (u TaskProjectIDUpdater) ToExpression() expression.UpdateBuilder {
	return expression.Set(expression.Name(projectIDField), expression.Value(u.ProjectID))
}
