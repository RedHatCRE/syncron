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
	"go.mongodb.org/mongo-driver/bson"
)

type Count struct {
	AllReportsWithVersion int32
	AllUniqueDeployments  int32
	AllReports            int64
}

type Calculations struct {
	Deployments int32
	Percentage  float32
}

type Deployment struct {
	ID            string     `json:"_id,omitempty"`
	AccountNumber string     `json:"account_number,omitempty"`
	CinderConf    CinderConf `json:"cinder_conf,omitempty"`
	ManilaConf    ManilaConf `json:"manila_conf,omitempty"`
	DeploymentID  string     `json:"deployment_id,omitempty"`
	Versions      Vhistory   `json:"versions_history,omitempty"`
	GlanceConf    GlanceConf `json:"glance_api_conf,omitempty"`
}

type Vhistory struct {
	OSPv OspVer `json:"osp_ver,omitempty"`
}

type OspVer struct {
	OSP10   string `json:"10,omitempty"`
	OSP13   string `json:"13,omitempty"`
	OSP161  string `json:"16.1,omitempty"`
	OSP16_2 string `json:"16.2,omitempty"`
	OSP17   string `json:"17,omitempty"`
}

type OctaviaWorker struct {
	Time   string
	Number int
}
type GlanceConf struct {
	Default DefaultGlance `json:"default,omitempty"`
}
type DefaultGlance struct {
	GlanceDrivers string `json:"enabled_backends,omitempty"`
}
type CinderBackEnd struct {
	Driver string `json:"volume_driver,omitempty"`
}

type CinderConf struct {
	DatabaseKey       DBKey         `json:"database,omitempty"`
	Default           DefaultCinder `json:"default,omitempty"`
	KeyManager        KeyManagerkey `json:"key_manager,omitempty"`
	KeystoneAuthToken KeystoneAuth  `json:"keystone_authtoken,omitempty"`
	NovaKey           Nova          `json:"nova,omitempty"`
	TripleoNFS        TripleoNFSKey `json:"tripleo_nfs,omitempty"`
	UploadDate        string        `json:"up_data_date,omitempty"`
	Driver            []string      `json:"up_cinder_volume_drivers,omitempty"`
	BackendName       CinderBackEnd
}

type DBKey struct {
	Connection string `json:"connection,omitempty"`
	MaxRetries string `json:"db_max_retries,omitempty"`
}

type DefaultCinder struct {
	VolumeType         string `json:"default_volume_type,omitempty"`
	EnabledV3api       string `json:"enable_v3_api,omitempty"`
	EnabledBackends    string `json:"enabled_backends,omitempty"`
	OSApiVolumeWorkers string `json:"osapi_volume_workers,omitempty"`
	StoraAvailZone     string `json:"storage_availability_zone,omitempty"`
}

type KeyManagerkey struct {
	Backend string `json:"backend,omitempty"`
}

type KeystoneAuth struct {
	AuthType          string
	AuthUri           string
	ProjectDomainName string
	ProjectName       string
	UserDomainName    string
	Username          string
}

type Nova struct {
	AuthType          string
	AuthUrl           string
	Interface         string
	ProjectDomainName string
	ProjectName       string
	UserDomainName    string
	Username          string
}

type TripleoNFSKey struct {
	BackEndHost              string
	NasSecureFileOperations  bool
	NasSecureFilePermissions bool
	BFSMountOptions          []string
	NFSnapshotSupport        bool
	VolumeBackendName        string
	VolumeDriver             string
}

type CinderDriversKey struct {
	Drivers []string
}

type ManilaConf struct {
	Default DefaultManila `json:"default,omitempty"`
	Cephfs  Cephfs        `json:"cephfs,omitempty"`
}

type DefaultManila struct {
	Backends string `json:"enabled_share_backends,omitempty"`
}

type Cephfs struct {
	Driver string `json:"share_driver,omitempty"`
}

var Manila Calculations
var MyFilter bson.A
var Numbers Count
