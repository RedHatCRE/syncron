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

type StoragePipelines struct {
	GlanceConf       bson.A
	NovaConf         bson.A
	ManilaEnabled    bson.A
	Cephfs_enabled   bson.A
	TripleoCeph      bson.A
	CephEnabled      bson.A
	DirectorDeployed bson.A
	NFSManila        bson.A
	HciDeploy        bson.A
}

var StoragePipes StoragePipelines
var Ga []bson.M

func InitStoragePipes() StoragePipelines {
	Ga = []bson.M{
		{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
		{"$match": bson.M{"glance_ips": bson.M{"$exists": true}}},
	}
	StoragePipes := StoragePipelines{
		GlanceConf: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"glance_ips": bson.M{"$exists": true}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		NovaConf: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"nova_conf": bson.M{"$exists": true}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		ManilaEnabled: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"$or": []bson.M{
				{"manila_conf": bson.M{"$exists": true}},
				{"ps.manila_manage": bson.M{"$exists": true}}}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		Cephfs_enabled: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{
				"manila_conf.default.enabled_share_backends": bson.M{"$in": []interface{}{"cephfs"}},
			}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		TripleoCeph: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"cinder_conf.default.enabled_backends": "tripleo_ceph"}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		CephEnabled: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"$or": []bson.M{
				{"glance_api_conf.glance_store.default_store": "rbd"},
				{"glance_api_conf.up_glance_default_store": "rbd"},
				{"glance_api_conf.default.enabled_backends": "default_backend:rbd"},
				{"ps.ceph-osd": bson.M{"$exists": true}},
				{"cinder_conf.up_cinder_volume_drivers": bson.M{"$in": []string{"rdb"}}}}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		DirectorDeployed: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"ps.ceph-mon": bson.M{"$exists": true}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		NFSManila: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"cinder_conf.up_cinder_volume_drivers": bson.M{"$in": []string{"nfs"}}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
		HciDeploy: bson.A{
			bson.M{"$match": bson.M{"versions_history.osp_ver": bson.M{"$exists": true}}},
			bson.M{"$match": bson.M{"ps.ceph-osd": bson.M{"$exists": true}}},
			bson.M{"$group": bson.M{"_id": "$deployment_id"}},
			bson.M{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		},
	}
	return StoragePipes
}
