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

package pipes

import "go.mongodb.org/mongo-driver/bson"

type CommonPipelines struct {
	TotalWithVersion  bson.A
	UniqueDeployments bson.A
	GatherVersions    bson.A
}

var CommonPipes CommonPipelines

func InitCommonPipes() CommonPipelines {
	StoragePipes := CommonPipelines{

		TotalWithVersion: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},

			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},

		UniqueDeployments: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": "deployment_id", "count": bson.M{"$sum": 1}}},
		},
		GatherVersions: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{
				"$project": bson.M{"firstkey": bson.M{"$arrayElemAt": []interface{}{bson.M{"$objectToArray": "$versions_history.osp_ver"}, 0}},
					"deployment_id": "$deployment_id",
				}},
			bson.M{"$group": bson.M{"_id": "$deployment_id", "uniqueValues": bson.M{"$first": "$firstkey.k"}}},
			bson.M{"$group": bson.M{"_id": "$uniqueValues", "count": bson.M{"$sum": 1}}},
		},
	}
	return StoragePipes
}
