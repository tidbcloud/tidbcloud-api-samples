// Copyright(C) 2022 PingCAP. All Rights Reserved.

package main

import (
	"fmt"
	"os"
)

const (
	Host = "https://api.tidbcloud.com"
)

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

// createDevCluster create a cluster in the given project
func createDevCluster(projectID uint64) (*CreateClusterResp, error) {
	var (
		url    = fmt.Sprintf("%s/api/v1beta/projects/%d/clusters", Host, projectID)
		result CreateClusterResp
	)

	payload := CreateClusterReq{
		Name:          "tidbcloud-sample-1", // NOTE change to your cluster name
		ClusterType:   "DEVELOPER",
		CloudProvider: "AWS",
		Region:        "us-east-1",
		Config: ClusterConfig{
			RootPassword: "your secret password", // NOTE change to your cluster password, we generate a random password here
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
	projects, err := getAllProjects()
	if err != nil {
		fmt.Printf("Failed to get all projects: %s\n", err)
		return
	}

	// Quit if we don't have any project
	if len(projects) < 1 {
		fmt.Printf("Failed to find any project in current organization\n")
		return
	}

	// Step 2: select our first project, and then create a cluster
	project := projects[0]
	fmt.Printf("Step 2: create developer cluster in project %d\n", project.ID)
	cluster, err := createDevCluster(project.ID)
	if err != nil {
		fmt.Printf("Failed to create developer cluster: %s\n", err)
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

	fmt.Printf("\nYou have created a developer cluster(don't forget to delete it).\n")
}
