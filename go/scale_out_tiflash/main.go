// Copyright(C) 2022 PingCAP. All Rights Reserved.

package main

import (
	"fmt"
	"os"
	"strconv"
)

const (
	Host = "https://api.tidbcloud.com"
)

// getSpecifications returns all the available specifications
func getSpecifications() (*GetSpecificationsResp, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/clusters/provider/regions", Host)
		result GetSpecificationsResp
	)

	_, err := doGET(url, nil, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getSpecClusterUsing(specifications *GetSpecificationsResp, cluster *GetClusterResp) (*Specification, error) {
	for _, i := range specifications.Items {
		if i.ClusterType == cluster.ClusterType && i.CloudProvider == cluster.CloudProvider && i.Region == cluster.Region {
			return &i, nil
		}
	}

	return nil, fmt.Errorf("No specification found")
}

// addTiFlashToCluster add TiFlash to a existing cluster
// If the vCPUs of TiDB or TiKV component is 2 or 4, then the cluster does not support TiFlash.
func addTiFlashToCluster(cluster *GetClusterResp, projectID uint64, spec *Specification) error {
	var (
		url              = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d", Host, projectID, cluster.ID)
		tiFlashComponent = cluster.Config.Components.TiFlash
	)

	if tiFlashComponent.NodeSize == "" {
		if len(spec.Tiflash) == 0 {
			return fmt.Errorf("No specification found")
		}

		// add TiFlash for cluster
		tiFlashSpec := spec.Tiflash[0]
		tiFlashComponent.NodeSize = tiFlashSpec.NodeSize
		tiFlashComponent.StorageSizeGib = tiFlashSpec.StorageSizeGibRange.Min
		tiFlashComponent.NodeQuantity = tiFlashSpec.NodeQuantityRange.Step
	} else {
		tiFlashComponent.NodeQuantity += 1 // add 1 TiFlash node it already has TiFlash
	}

	payload := ScaleClusterReq{
		Config: ClusterConfig{
			Components: Components{
				TiDB: ComponentTiDB{
					NodeSize:       cluster.Config.Components.TiDB.NodeSize,
					StorageSizeGib: cluster.Config.Components.TiDB.StorageSizeGib,
					NodeQuantity:   cluster.Config.Components.TiDB.NodeQuantity,
				},
				TiKV: ComponentTiKV{
					NodeSize:       cluster.Config.Components.TiKV.NodeSize,
					StorageSizeGib: cluster.Config.Components.TiKV.StorageSizeGib,
					NodeQuantity:   cluster.Config.Components.TiKV.NodeQuantity,
				},
				TiFlash: tiFlashComponent,
			},
		},
	}

	_, err := doPATCH(url, payload, nil)
	return err
}

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

	// Step 1: scale cluster, add one TiFlash node
	specs, err := getSpecifications()
	if err != nil {
		fmt.Printf("Failed to get specifications: %s\n", err)
		return
	}

	cluster, err := getClusterByID(projectID, clusterID)
	if err != nil {
		fmt.Printf("Failed to get cluster %d: %s", clusterID, err)
		return
	}
	if cluster.Status.ClusterStatus != ClusterAvailable {
		fmt.Printf("Bad cluster status: %s", cluster.Status.ClusterStatus)
		return
	}

	spec, err := getSpecClusterUsing(specs, cluster)
	if err != nil {
		fmt.Printf("Failed to get dedicated specification: %s\n", err)
		return
	}

	fmt.Printf("Step 1: add TiFlash to cluster %d\n", clusterID)
	err = addTiFlashToCluster(cluster, projectID, spec)
	if err != nil {
		fmt.Printf("Failed to add TiFlash to cluster %d: %s\n", clusterID, err)
		return
	}

	// Step 2: check progress
	fmt.Printf("Step 2: wait cluster to get ready\n")
	cluster, err = getClusterByID(projectID, clusterID)
	if err != nil {
		fmt.Printf("Failed to get cluster status: %s\n", err)
		return
	}
	fmt.Printf("Cluster status: %+v\n", cluster)
	//waitFor("cluster", func() bool {
	//cluster, err := getClusterByID(projectID, clusterID)
	//if err != nil {
	//fmt.Printf("Failed to query cluster detail by id %d: %s\n", clusterID, err)
	//return false
	//}

	//return cluster.Status.ClusterStatus == ClusterAvailable
	//})

	fmt.Printf("You have scaled the cluster and add TiFlash\n")
}
