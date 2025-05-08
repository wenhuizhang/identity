// Copyright 2025  AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
)

func Match(query interface{}) bson.D {
	return bson.D{{Key: "$match", Value: query}}
}

func Project(value interface{}) bson.D {
	return bson.D{{Key: "$project", Value: value}}
}

func Unwind(value string) bson.D {
	return bson.D{
		{
			Key: "$unwind", Value: value,
		},
	}
}

func ReplaceRoot(value string) bson.D {
	return bson.D{
		{
			Key:   "$replaceRoot",
			Value: bson.D{{Key: "newRoot", Value: value}},
		},
	}
}
