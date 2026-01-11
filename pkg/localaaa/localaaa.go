// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package localaaa

import (
	"context"

	"github.com/itential/ipctl/internal/logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	localAaaDatabase           = "LocalAAA"
	localAaaAccountsCollection = "accounts"
	localAaaGroupsCollection   = "groups"
)

type LocalAAAService struct {
	accounts *mongo.Collection
	groups   *mongo.Collection
}

func NewLocalAAAService(uri string) LocalAAAService {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		logging.Fatal(err, "failed to connect to db")
	}

	return LocalAAAService{
		accounts: client.Database(localAaaDatabase).Collection(localAaaAccountsCollection),
		groups:   client.Database(localAaaDatabase).Collection(localAaaGroupsCollection),
	}
}
