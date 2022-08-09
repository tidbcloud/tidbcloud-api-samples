# Copyright(C) 2022 PingCAP. All Rights Reserved.
"""
Purpose

Shows how to create a TiDB cluster and get cluster's detail.
"""
import json
import os
import time

import requests
from requests.auth import HTTPDigestAuth

# Basic config
HOST = "https://api.tidbcloud.com"


class CreateDedicatedCluster:
    def __init__(self):
        """
        Generate authorization by public key and private key (https://tidbcloud.com/console/clusters)
        """
        self.digest_auth = _authorization()

    def get_all_projects(self):
        """
        Get all projects
        :return: Projects detail
        """
        url = f"{HOST}/api/v1beta/projects"
        resp = requests.get(url=url, auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        return _response(resp)

    def get_provider_regions_specifications(self):
        """
        Get cloud providers, regions and available specifications.
        :return: List the cloud providers, regions and available specifications.
        """
        url = f"{HOST}/api/v1beta/clusters/provider/regions"
        resp = requests.get(url=url, auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        return _response(resp)

    def create_dedicated_cluster(self, project_id: str, dedicated_config: dict) -> dict:
        """
        Create a dedicated cluster in your specified project
        :param project_id: The project id
        :param dedicated_config: The dedicated config
        :return: Dedicated cluster id
        """
        try:
            region = dedicated_config["region"]
            cloud_provider = dedicated_config["cloud_provider"]
            tidb_size = dedicated_config["tidb"][0]["node_size"]
            tidb_quantity_range = dedicated_config["tidb"][0]["node_quantity_range"]["min"]
            tikv_size = dedicated_config["tikv"][0]["node_size"]
            tikv_quantity_range = dedicated_config["tikv"][0]["node_quantity_range"]["min"]
            tikv_storage_size_gib_range = dedicated_config["tikv"][0]["storage_size_gib_range"]["min"]
        except (KeyError, IndexError) as e:
            print(f"cloud provider or region or available specifications not found! exception: {e}")
            raise
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters"
        ts = int(time.time())
        data_config = \
            {
                "name": f"tidbcloud-sample-{ts}",
                "cluster_type": "DEDICATED",
                "cloud_provider": f"{cloud_provider}",
                "region": f"{region}",
                "config":
                    {
                        "root_password": "input_your_password",
                        "port": 4000,
                        "components":
                            {
                                "tidb":
                                    {
                                        "node_size": f"{tidb_size}",
                                        "node_quantity": f"{tidb_quantity_range}"
                                    },
                                "tikv":
                                    {
                                        "node_size": f"{tikv_size}",
                                        "storage_size_gib": f"{tikv_storage_size_gib_range}",
                                        "node_quantity": f"{tikv_quantity_range}"
                                    }
                            },
                        "ip_access_list":
                            [
                                {
                                    "cidr": "0.0.0.0/0",
                                    "description": "Allow Access from Anywhere."
                                }
                            ]

                    }
            }
        data_config_json = json.dumps(data_config)
        resp = requests.post(url=url,
                             auth=self.digest_auth,
                             data=data_config_json)
        print(f"Method: {resp.request.method}, Request: {url}, Payload: {data_config_json}")
        return _response(resp)

    def get_cluster_by_id(self, project_id: str, cluster_id: str) -> dict:
        """
        Get the cluster detail.
        You will get `connection_strings` from the response after the cluster's status is`AVAILABLE`.
        Then, you can connect to TiDB using the default user, host, and port in `connection_strings`
        :param project_id: The project id
        :param cluster_id: The cluster id
        :return: The cluster detail
        """
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters/{cluster_id}"
        resp = requests.get(url=url,
                            auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        return _response(resp)

    def delete_cluster(self, project_id: str, cluster_id: str) -> dict:
        """
        Delete the cluster
        :param project_id: The project id
        :param cluster_id: The cluster id
        :return: Result for deletion
        """
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters/{cluster_id}"
        resp = requests.delete(url=url,
                               auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        return _response(resp)


def _authorization():
    """
    Generate authorization by public key and private key (https://tidbcloud.com/console/clusters)
    :return: Digest auth
    """
    public_key = os.environ.get("TIDBCLOUD_PUBLIC_KEY", None)
    private_key = os.environ.get("TIDBCLOUD_PRIVATE_KEY", None)
    if public_key is None or private_key is None:
        print("TIDBCLOUD_PUBLIC_KEY or TIDBCLOUD_PRIVATE_KEY is None, you should set TIDBCLOUD_PUBLIC_KEY and "
              "TIDBCLOUD_PRIVATE_KEY firstly.")
        raise Exception("TIDBCLOUD_PUBLIC_KEY or TIDBCLOUD_PRIVATE_KEY not set.")
    return HTTPDigestAuth(public_key, private_key)


def _response(resp: requests.models.Response) -> dict:
    """
    Response from open api
    :param resp: Result from the API
    :return: Format response
    """
    if resp.status_code != 200:
        print(f"request invalid, code : {resp.status_code}, message : {resp.text}")
        raise Exception(f"request invalid, code : {resp.status_code}, message : {resp.text}")
    print(f"response : {resp.text}")
    return resp.json()


def _get_dedicated_config(provider_region_specifications: dict) -> dict:
    """
    Get cloud providers, regions and available specifications of dedicated
    :param provider_region_specifications: Result of  cloud providers, regions and available specifications
    :return: Available dedicated config
    """
    try:
        items = provider_region_specifications["items"]
        for item in items:
            if item["cluster_type"] == "DEDICATED":
                return item
    except KeyError:
        print("dedicate config not found!")
        raise
    raise Exception("dedicate config not found!")


def usage_demo():
    print("-" * 88)
    print("Welcome to the TiDB Cloud API samples!")
    print("-" * 88)

    create_cluster = CreateDedicatedCluster()

    print("1. Get all projects. ")
    projects = create_cluster.get_all_projects()
    try:
        items = projects["items"]
        if len(items) > 0:
            sample_project_id = items[0]["id"]
    except KeyError:
        print("project id not found!")
        raise
    print()

    print("2. Get the cloud providers, regions and available specifications")
    provider_region_specifications = create_cluster.get_provider_regions_specifications()
    dedicated_config = _get_dedicated_config(provider_region_specifications)
    print()

    print(f"3. Create a dedicated cluster in your specified project ( project id : {sample_project_id} ). ")
    cluster = create_cluster.create_dedicated_cluster(sample_project_id, dedicated_config)
    try:
        sample_cluster_id = cluster["id"]
    except KeyError:
        print("cluster id not found!")
        raise
    print()

    print(f"4. Get the new cluster ( cluster id : {sample_cluster_id} ) detail. ")
    create_cluster.get_cluster_by_id(sample_project_id, sample_cluster_id)
    print()

    # print("If necessary , delete the cluster.")
    # create_cluster.delete_cluster(sample_project_id, sample_cluster_id)

    print("Thanks for watching! ")
    print("-" * 88)


if __name__ == "__main__":
    usage_demo()
