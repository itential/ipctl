// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package localaaa

import (
	"context"

	"github.com/itential/ipctl/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

type Group struct {
	Name string `json:"name" bson:"name"`
}

func NewGroup(name string) Group {
	return Group{Name: name}
}

func (svc LocalAAAService) GetGroups() ([]Group, error) {
	logger.Trace()

	cur, err := svc.groups.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var documents []bson.M
	if err := cur.All(context.TODO(), &documents); err != nil {
		return nil, err
	}

	var res []Group
	for _, ele := range documents {
		b, err := bson.Marshal(ele)
		if err != nil {
			return nil, err
		}

		grp := new(Group)

		if err := bson.Unmarshal(b, &grp); err != nil {
			return nil, err
		}

		res = append(res, *grp)
	}

	return res, nil
}

func (svc LocalAAAService) CreateGroup(in Group) error {
	logger.Trace()
	_, err := svc.groups.InsertOne(context.TODO(), in)
	return err

}

func (svc LocalAAAService) DeleteGroup(name string) error {
	logger.Trace()

	_, err := svc.groups.DeleteOne(
		context.TODO(),
		bson.D{{Key: "name", Value: name}},
	)

	return err
}
