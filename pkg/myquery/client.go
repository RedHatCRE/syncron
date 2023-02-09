// Copyright 2022 Red Hat, Inc.
// All Rights Reserved.
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package myquery

import (
	"context"

	"github.com/redhatcre/syncron/configuration"
	"github.com/redhatcre/syncron/pkg/cli"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MyDataBase struct {
	MyClient *mongo.Client
	Coll     *mongo.Collection
}

// This function creates a client for accessing a MongoDB database.
func (m *MyDataBase) CreateClient() *MyDataBase {

	GetUserInfo()
	c := configuration.Configuration{}
	c.GetConfiguration()
	uri := c.MongoConf.Uri

	credential := options.Credential{
		AuthMechanism: "PLAIN",
		Username:      User.Name,
		Password:      User.Pass,
	}

	s := cli.IniSpinner(" Creating client")
	var err error
	m.MyClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri).SetAuth(credential))
	if err != nil {
		s.Stop()
		logrus.Fatal("Error creating client", err)
	}
	if err = m.MyClient.Ping(context.Background(), readpref.Primary()); err != nil {
		s.Stop()
		logrus.Fatal("Error pinging mongodb client.")
	} else {
		s.Stop()
		logrus.Info("Success setting up your Client")
	}
	m.Coll = m.MyClient.Database(c.MongoConf.Database).Collection(c.MongoConf.Collection)
	return m
}
