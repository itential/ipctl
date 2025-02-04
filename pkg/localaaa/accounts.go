// Copyright 2024 Itential Inc. All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential

package localaaa

import (
	"context"
	"errors"

	"github.com/itential/ipctl/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	//Id           bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username     string   `json:"username" bson:"username"`
	ActiveTenant string   `json:"activeTenant" bson:"activeTenant"`
	FirstName    string   `json:"firstname" bson:"firstname"`
	Groups       []string `json:"groups" bson:"groups"`
	Password     string   `json:"password" bson:"password"`
	Tenants      []string `json:"tenants" bson:"tenants"`
}

func NewAccount(username, password string) Account {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		logger.Fatal(err, "failed to encrypt value")
	}

	return Account{
		Username:     username,
		Password:     string(b),
		ActiveTenant: "*",
		Tenants:      []string{},
		Groups:       []string{},
	}
}

func (svc LocalAAAService) GetAccounts() ([]Account, error) {
	logger.Trace()

	cur, err := svc.accounts.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var documents []bson.M
	if err := cur.All(context.TODO(), &documents); err != nil {
		return nil, err
	}

	var res []Account
	for _, ele := range documents {
		b, err := bson.Marshal(ele)
		if err != nil {
			return nil, err
		}

		user := new(Account)

		if err := bson.Unmarshal(b, &user); err != nil {
			return nil, err
		}

		res = append(res, *user)
	}

	return res, nil
}

func (svc LocalAAAService) Find(username string) (*Account, error) {
	logger.Trace()

	sr := svc.accounts.FindOne(context.TODO(), bson.D{
		{Key: "username", Value: username},
	})

	var res *Account

	if sr.Err() == nil {
		if err := sr.Decode(&res); err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (svc LocalAAAService) CreateAccount(in Account) error {
	logger.Trace()

	res, err := svc.Find(in.Username)
	if err != nil {
		return err
	}

	if res != nil {
		return errors.New("username already exists")
	}

	_, err = svc.accounts.InsertOne(context.TODO(), in)
	if err != nil {
		return err
	}

	return nil
}

func (svc LocalAAAService) DeleteAccount(username string) error {
	logger.Trace()

	_, err := svc.accounts.DeleteOne(
		context.TODO(),
		bson.D{{Key: "username", Value: username}},
	)

	return err
}
