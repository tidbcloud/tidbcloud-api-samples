# -*- coding: utf-8 -*-
# Copyright(C) 2022 PingCAP. All Rights Reserved.
"""
Purpose

Shows how to scale-out one Tiflash node for existing cluster.
"""
import json
import os

import requests
from requests.auth import HTTPDigestAuth

# Basic config
HOST = "https://api.tidbcloud.com"


class ScaleOutTiFlash:
    def __init__(self):
        """
        Generate authorization by public key and private key (https://tidbcloud.com/console/clusters)
        """
        self.digest_auth = _authorization()

    def get_dedicated_provider_regions_specifications(self):
        """
        Get dedicated cloud providers, regions and available specifications.
        :return: List the cloud providers, regions and available specifications.
        """
        url = f"{HOST}/api/v1beta/clusters/provider/regions"
        resp = requests.get(url=url, auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        resp_body = _response(resp)
        try:
            items = resp_body["items"]
            for item in items:
                if item["cluster_type"] == "DEDICATED":
                    return item
        except KeyError:
            print("dedicate config not found!")
            raise
        raise Exception("dedicate config not found!")

    def modify_cluster(self, project_id: str, cluster_id: str, dedicated_config: dict) -> dict:
        """
        Add one TiFlash node for specified cluster
        :param project_id: The project id
        :param cluster_id: The cluster id
        :param dedicated_config: The dedicated config
        :return: If success,return None. Else, return message
        """
        try:
            provider_regions_specifications = self.get_dedicated_provider_regions_specifications()
            tidb_quantity_range = dedicated_config["config"]["components"]["tidb"]["node_quantity"]
            tikv_quantity_range = dedicated_config["config"]["components"]["tikv"]["node_quantity"]
            if dedicated_config["config"]["components"]["tiflash"] is not None:
                tiflash_size = dedicated_config["config"]["components"]["tiflash"]["node_size"]
                tiflash_storage_size_gib = dedicated_config["config"]["components"]["tiflash"]["storage_size_gib"]
                tiflash_node_quantity = dedicated_config["config"]["components"]["tiflash"]["node_quantity"] + 1
            else:
                tiflash_size = provider_regions_specifications["tiflash"][0]["node_size"]
                tiflash_storage_size_gib = provider_regions_specifications["tiflash"][0]["storage_size_gib_range"][
                    "min"]
                tiflash_node_quantity = provider_regions_specifications["tiflash"][0]["node_quantity_range"]["step"]
        except (KeyError, IndexError) as e:
            print(f"cloud provider or region or available specifications not found! exception: {e}")
            raise
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters/{cluster_id}"
        data_add_tiflash = \
            {
                "config":
                    {
                        "components":
                            {
                                "tidb":
                                    {
                                        "node_quantity": f"{tidb_quantity_range}"
                                    },
                                "tikv":
                                    {
                                        "node_quantity": f"{tikv_quantity_range}"
                                    },
                                "tiflash":
                                    {
                                        "node_quantity": f"{tiflash_node_quantity}",
                                        "node_size": f"{tiflash_size}",
                                        "storage_size_gib": f"{tiflash_storage_size_gib}"
                                    }
                            }
                    }
            }
        data_add_tiflash_json = json.dumps(data_add_tiflash)
        resp = requests.patch(url=url,
                              auth=self.digest_auth,
                              data=data_add_tiflash_json)
        print(f"Method: {resp.request.method}, Request: {url}, Payload:{data_add_tiflash_json}")
        return _response(resp)

    def get_cluster_by_id(self, project_id: str, cluster_id: str) -> dict:
        """
        Get the cluster detail
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


def usage_demo():
    print("-" * 88)
    print("Welcome to the TiDB Cloud API samples!")
    print("-" * 88)

    scale_out_tiflash = ScaleOutTiFlash()
    dedicated_project_id = os.environ.get("DEDICATED_PROJECT_ID", None)
    dedicated_cluster_id = os.environ.get("DEDICATED_CLUSTER_ID", None)
    if dedicated_project_id is None or dedicated_cluster_id is None:
        print(
            "DEDICATED_PROJECT_ID or DEDICATED_CLUSTER_ID is None, you should set DEDICATED_PROJECT_ID and "
            "DEDICATED_CLUSTER_ID firstly.")
        raise Exception("DEDICATED_PROJECT_ID or DEDICATED_CLUSTER_ID not set!")

    print("1. Add one TiFlash node for specified cluster.")
    dedicated_cluster_info = scale_out_tiflash.get_cluster_by_id(dedicated_project_id, dedicated_cluster_id)
    scale_out_tiflash.modify_cluster(dedicated_project_id, dedicated_cluster_id, dedicated_cluster_info)
    print()

    print("2. View the scale-out progress.")
    scale_out_tiflash.get_cluster_by_id(dedicated_project_id, dedicated_cluster_id)
    print()

    print("Thanks for watching! ")
    print("-" * 88)


if __name__ == "__main__":
    usage_demo()
