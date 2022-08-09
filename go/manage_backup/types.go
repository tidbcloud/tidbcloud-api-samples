// Copyright(C) 2022 PingCAP. All Rights Reserved.

package main

import (
	"time"
)

const (
	// Time format
	DateFormat = "2006-01-02"

	// Backup status
	BackupSuccess = "SUCCESS"

	// Cluster status
	ClusterAvailable = "AVAILABLE"
)

type Project struct {
	ID              uint64 `json:"id,string"`
	OrgID           uint64 `json:"org_id,string"`
	Name            string `json:"name"`
	ClusterCount    int64  `json:"cluster_count"`
	UserCount       int64  `json:"user_count"`
	CreateTimestamp int64  `json:"create_timestamp,string"`
}

type ConnectionString struct {
	Standard   string `json:"standard"`
	VpcPeering string `json:"vpc_peering"`
}

type IPAccess struct {
	CIDR        string `json:"cidr"`
	Description string `json:"description"`
}

type ComponentTiDB struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type ComponentTiKV struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type ComponentTiFlash struct {
	NodeSize       string `json:"node_size"`
	StorageSizeGib int32  `json:"storage_size_gib"`
	NodeQuantity   int32  `json:"node_quantity"`
}

type Components struct {
	TiDB    *ComponentTiDB    `json:"tidb,omitempty"`
	TiKV    *ComponentTiKV    `json:"tikv,omitempty"`
	TiFlash *ComponentTiFlash `json:"tiflash,omitempty"`
}

type ClusterConfig struct {
	RootPassword string     `json:"root_password"`
	Port         int32      `json:"port"`
	Components   Components `json:"components"`
	IPAccessList []IPAccess `json:"ip_access_list"`
}

type ClusterStatus struct {
	TidbVersion   string `json:"tidb_version"`
	ClusterStatus string `json:"cluster_status"`
}

type GetClusterBackupResp struct {
	ID              uint64    `json:"id,string"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	Type            string    `json:"type"`
	CreateTimestamp time.Time `json:"create_timestamp"`
	Size            uint64    `json:"size,string"`
	Status          string    `json:"status"`
}

type CreateClusterReq struct {
	Name          string        `json:"name"`
	ClusterType   string        `json:"cluster_type"`
	CloudProvider string        `json:"cloud_provider"`
	Region        string        `json:"region"`
	Config        ClusterConfig `json:"config"`
}

type CreateClusterResp struct {
	ClusterID uint64 `json:"id,string"`
	Message   string `json:"message"`
}

type GetAllProjectsResp struct {
	Items []Project `json:"items"`
	Total int64     `json:"total"`
}

type GetClusterResp struct {
	ID                uint64           `json:"id,string"`
	ProjectID         uint64           `json:"project_id,string"`
	Name              string           `json:"name"`
	Port              int32            `json:"port"`
	TiDBVersion       string           `json:"tidb_version"`
	ClusterType       string           `json:"cluster_type"`
	CloudProvider     string           `json:"cloud_provider"`
	Region            string           `json:"region"`
	Status            ClusterStatus    `json:"status"`
	CreateTimestamp   string           `json:"create_timestamp"`
	Config            ClusterConfig    `json:"config"`
	ConnectionStrings ConnectionString `json:"connection_strings"`
}

type CreateclusterBackupReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateclusterBackupResp struct {
	ID uint64 `json:"id,string"`
}

type RestoreClusterReq struct {
	BackupID uint64        `json:"backup_id,string"`
	Name     string        `json:"name"`
	Config   ClusterConfig `json:"config"`
}

type RestoreClusterResp struct {
	ID        uint64 `json:"id,string"`
	ClusterID uint64 `json:"cluster_id,string"`
}
