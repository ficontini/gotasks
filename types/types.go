package types

import (
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minTitleLen       = 5
	minDescriptionLen = 5
)

type ID string

func CreateIDFromObjectID(oid primitive.ObjectID) ID {
	return ID(oid.Hex())
}
func CreateIDFromInt64(id int64) ID {
	return ID(strconv.Itoa(int(id)))
}
func (id ID) Int() (int64, error) {
	return strconv.ParseInt(string(id), 10, 64)
}

func (id ID) String() string {
	return string(id)
}
func (id ID) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id.String())
}

func (id *ID) Scan(value interface{}) error {
	intValue, ok := value.(int64)
	if !ok {
		return fmt.Errorf("unexpected type for ID %T", value)
	}
	*id = CreateIDFromInt64(intValue)
	return nil
}
