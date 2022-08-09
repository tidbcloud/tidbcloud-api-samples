// Copyright(C) 2022 PingCAP. All Rights Reserved.

package main

import (
	"fmt"
	"os"
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

func getDedicatedSpec(specifications *GetSpecificationsResp) (*Specification, error) {
	for _, i := range specifications.Items {
		if i.ClusterType == "DEDICATED" {
			return &i, nil
		}
	}

	return nil, fmt.Errorf("No specification found")
}

// getAllProjects list all projects in current organization
func getAllProjects() ([]Project, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects", Host)
		result GetAllProjectsResp
	)

	_, err := doGET(url, nil, &result)
	if err != nil {
		return nil, err
	}

	return result.Items, nil
}

// createDedicatedCluster create a cluster in the given project
func createDedicatedCluster(projectID uint64, spec *Specification) (*CreateClusterResp, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters", Host, projectID)
		result CreateClusterResp
	)

	// We have check the boundary in main function
	tidbSpec := spec.Tidb[0]
	tikvSpec := spec.Tikv[0]

	payload := CreateClusterReq{
		Name:          "tidbcloud-sample-1", // NOTE change to your cluster name
		ClusterType:   spec.ClusterType,
		CloudProvider: spec.CloudProvider,
		Region:        spec.Region,
		Config: ClusterConfig{
			RootPassword: "your secret password", // NOTE change to your cluster password, we generate a random password here
			Port:         4000,
			Components: Components{
				TiDB: &ComponentTiDB{
					NodeSize:     tidbSpec.NodeSize,
					NodeQuantity: tidbSpec.NodeQuantityRange.Min,
				},
				TiKV: &ComponentTiKV{
					NodeSize:       tikvSpec.NodeSize,
					StorageSizeGib: tikvSpec.StorageSizeGibRange.Min,
					NodeQuantity:   tikvSpec.NodeQuantityRange.Min,
				},
			},
			IPAccessList: []IPAccess{
				{
					CIDR:        "0.0.0.0/0",
					Description: "Allow Access from Anywhere.",
				},
			},
		},
	}

	_, err := doPOST(url, payload, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
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

// deleteClusterByID delete a cluster
func deleteClusterByID(projectID, clusterID uint64) error {
	url := fmt.Sprintf("%s/api/v1beta/projects/%d/clusters/%d", Host, projectID, clusterID)
	_, err := doDELETE(url, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Printf("\nWelcome to the TiDB Cloud demo!\n\n")

	// Step 0: we need to initialize http client first
	var (
		publicKey  = os.Getenv("TIDBCLOUD_PUBLIC_KEY")
		privateKey = os.Getenv("TIDBCLOUD_PRIVATE_KEY")
	)
	if publicKey == "" || privateKey == "" {
		fmt.Printf("Please set TIDBCLOUD_PUBLIC_KEY(%s), TIDBCLOUD_PRIVATE_KEY(%s) in environment variable first\n", publicKey, privateKey)
		return
	}

	fmt.Printf("Step 0: initialize HTTP client\n")
	err := initClient(publicKey, privateKey)
	if err != nil {
		fmt.Printf("Failed to init HTTP client\n")
		return
	}

	// Step 1: query all the projects in current organization
	fmt.Printf("Step 1: get all projects\n")
	specs, err := getSpecifications()
	if err != nil {
		fmt.Printf("Failed to get specifications: %s\n", err)
		return
	}

	dedicatedSpec, err := getDedicatedSpec(specs)
	if err != nil {
		fmt.Printf("Failed to get dedicated specification: %s\n", err)
		return
	}

	projects, err := getAllProjects()
	if err != nil {
		fmt.Printf("Failed to get all projects: %s\n", err)
		return
	}

	// Quit if we don't have any project, or can't find a valid specification
	if len(projects) < 1 {
		fmt.Printf("Failed to find any project in current organization\n")
		return
	}

	if len(dedicatedSpec.Tidb) < 1 || len(dedicatedSpec.Tikv) < 1 {
		fmt.Printf("Invalid specification: no available TiDB(%d)/TiKV(%d) specifications\n", len(dedicatedSpec.Tidb), len(dedicatedSpec.Tikv))
		return
	}

	// Step 2: select our first project, and then create a cluster
	project := projects[0]
	fmt.Printf("Step 2: create dedicated cluster in project %d\n", project.ID)
	cluster, err := createDedicatedCluster(project.ID, dedicatedSpec)
	if err != nil {
		fmt.Printf("Failed to create dedicated cluster: %s\n", err)
		return
	}

	// Step 3: check the status of created cluster
	fmt.Printf("Step 3: get cluster by id %d\n", cluster.ClusterID)
	resp, err := getClusterByID(project.ID, cluster.ClusterID)
	if err != nil {
		fmt.Printf("Failed to get cluster detail of %d: %s\n", cluster.ClusterID, err)
		return
	}

	fmt.Printf("Cluster detail: %+v\n", resp)
	//You will get `connection_strings` from the response after the cluster's status is
	//`AVAILABLE`. Then, you can connect to TiDB using the default user, host, and port in `connection_strings`.
	//fmt.Printf("Connection string: %s", resp.ConnectionStrings.Standard)

	// Step 4: delete the newly created cluster
	//fmt.Printf("Step 4: delete cluster by id %d\n", cluster.ClusterID)
	//err = deleteClusterByID(project.ID, cluster.ClusterID)
	//if err != nil {
	//fmt.Printf("Failed to delete cluster %d: %s\n", cluster.ClusterID, err)
	//return
	//}

	fmt.Printf("\nYou have created a dedicated cluster(don't forget to delete it).\n")
}
