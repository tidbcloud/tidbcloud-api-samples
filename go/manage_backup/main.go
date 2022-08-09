// Copyright(C) 2022 PingCAP. All Rights Reserved.

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	Host = "https://api.tidbcloud.com"
)

// getClusterByID return detail status of given cluster
func getClusterByID(projectID, clusterID uint64) (*GetClusterResp, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d", Host, projectID, clusterID)
		result GetClusterResp
	)

	_, err := doGET(url, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// getBackupByID get backup detail by backup id
func getBackupByID(projectID, clusterID, backupID uint64) (*GetClusterBackupResp, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d/backups/%d", Host, projectID, clusterID, backupID)
		result GetClusterBackupResp
	)

	_, err := doGET(url, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// deleteBackupByID delete a backup by backup id
func deleteBackupByID(projectID, clusterID, backupID uint64) error {
	var (
		url = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d/backups/%d", Host, projectID, clusterID, backupID)
	)

	_, err := doDELETE(url, nil, nil)
	return err
}

// createBackupForCluster create backup for cluster
func createBackupForCluster(projectID, clusterID uint64) (uint64, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d/backups", Host, projectID, clusterID)
		result CreateclusterBackupResp
	)

	now := time.Now()
	payload := CreateclusterBackupReq{
		Name:        fmt.Sprintf("tidbcloud-backup-%s", now.Format(DateFormat)),
		Description: fmt.Sprintf("tidbcloud backup created for demo in %s", now.Format(DateFormat)),
	}
	_, err := doPOST(url, payload, &result)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

// deleteClusterByID delete cluster by cluster id
func deleteClusterByID(projectID, clusterID uint64) error {
	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d", Host, projectID, clusterID)
	_, err := doDELETE(url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

// restoreClusterByBackupID re-create a new cluster from backup
func restoreClusterByBackupID(cluster *GetClusterResp, projectID, backupID uint64) (*RestoreClusterResp, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects/%d/restores", Host, projectID)
		result RestoreClusterResp
	)

	payload := RestoreClusterReq{
		BackupID: backupID,
		Name:     "tidbcloud-sample-restore",
		Config: ClusterConfig{
			RootPassword: "your secret password", // NOTE change to your cluster password, we generate a random password here
			Port:         4000,
			Components: Components{
				TiDB:    cluster.Config.Components.TiDB,
				TiKV:    cluster.Config.Components.TiKV,
				TiFlash: cluster.Config.Components.TiFlash,
			},
		},
	}
	_, err := doPOST(url, payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func main() {
	fmt.Printf("\nWelcome to the TiDB Cloud demo!\n\n")

	// Step 0: we need to initialize http client first
	var (
		publicKey  = os.Getenv("TIDBCLOUD_PUBLIC_KEY")
		privateKey = os.Getenv("TIDBCLOUD_PRIVATE_KEY")
		projectID  uint64
		clusterID  uint64
	)
	if publicKey == "" || privateKey == "" {
		fmt.Printf("Please set TIDBCLOUD_PUBLIC_KEY(%s), TIDBCLOUD_PRIVATE_KEY(%s) in environment variable first\n", publicKey, privateKey)
		return
	}

	projectID, err := strconv.ParseUint(os.Getenv("DEDICATED_PROJECT_ID"), 10, 64)
	if err != nil {
		fmt.Printf("Please set DEDICATED_PROJECT_ID in environment variable first\n")
		return
	}

	clusterID, err = strconv.ParseUint(os.Getenv("DEDICATED_CLUSTER_ID"), 10, 64)
	if err != nil {
		fmt.Printf("Please set DEDICATED_CLUSTER_ID in environment variable first\n")
		return
	}

	fmt.Printf("Step 0: initialize HTTP client\n")
	err = initClient(publicKey, privateKey)
	if err != nil {
		fmt.Printf("Failed to init HTTP client\n")
		return
	}

	// Step 1: check the cluster
	fmt.Printf("Step 1: check cluster %d's status\n", clusterID)
	cluster, err := getClusterByID(projectID, clusterID)
	if err != nil {
		fmt.Printf("Failed to get cluster %d: %s", clusterID, err)
		return
	}
	if cluster.Status.ClusterStatus != ClusterAvailable {
		fmt.Printf("Bad cluster status: %s", cluster.Status.ClusterStatus)
		return
	}

	// Step 2: create a backup
	fmt.Printf("Step 2: create backup for cluster %d\n", clusterID)
	backupID, err := createBackupForCluster(projectID, clusterID)
	if err != nil {
		fmt.Printf("Failed to create cluster backup for cluster id %d: %s", clusterID, err)
		return
	}
	fmt.Printf("BackupID is %d\n", backupID)

	//waitFor("backup", func() bool {
	//backup, err := getBackupByID(projectID, clusterID, backupID)
	//if err != nil {
	//fmt.Printf("Failed to query cluster backup detail by id %d: %s", backupID, err)
	//return false
	//}

	//return backup.Status == BackupSuccess
	//})

	//// Step 3: restore from backup
	//fmt.Printf("Step 3: restore cluster from backup %d\n", backupID)
	//resp, err := restoreClusterByBackupID(cluster, projectID, backupID)
	//if err != nil {
	//fmt.Printf("Failed to restore cluster from backup id %d: %s", backupID, err)
	//return
	//}
	//cluster, err = getClusterByID(projectID, resp.ClusterID)
	//if err != nil {
	//fmt.Printf("Failed to get cluster status: %s\n", err)
	//return
	//}
	//fmt.Printf("Cluster status: %+v\n", cluster)

	//waitFor("cluster", func() bool {
	//cluster, err := getClusterByID(projectID, resp.ClusterID)
	//if err != nil {
	//fmt.Printf("Failed to query cluster detail by id %d: %s", resp.ClusterID, err)
	//return false
	//}

	//return cluster.Status.ClusterStatus == ClusterAvailable
	//})

	// Step 4: delete the backup
	//fmt.Printf("Step 4: delete backup by id %d\n", backupID)
	//if err = deleteBackupByID(projectID, clusterID, backupID); err != nil {
	//fmt.Printf("Failed to delete backup %d: %s", backupID, err)
	//return
	//}

	// Step 5: delete the newly created cluster
	//if err = deleteClusterByID(projectID, resp.ClusterID); err != nil {
	//fmt.Printf("Failed to delete cluster %d: %s", resp.ClusterID, err)
	//return
	//}

	fmt.Printf("You have created a backup and restore it to a new cluster(don't forget to delete it).\n")
}
