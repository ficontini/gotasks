package types

import "go.mongodb.org/mongo-driver/bson"

type Update interface {
	ToBSON() bson.M
}

type StatusUpdater struct {
	Enabled bool
}

func (e StatusUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"enabled": e.Enabled},
	}
}

type PasswordUpdater struct {
	EncryptedPassword string
}

func (e PasswordUpdater) ToBSON() bson.M {
	return bson.M{
		"$set": bson.M{"encryptedPassword": e.EncryptedPassword},
	}
}
