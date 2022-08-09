# Copyright(C) 2022 PingCAP. All Rights Reserved.
"""
Purpose

Shows how to  create a manual backup cyclically and restore last backup data to a new cluster.
"""
import datetime
import json
import os

import requests
from requests.auth import HTTPDigestAuth

# Basic config
HOST = "https://api.tidbcloud.com"


class ManageBackup:
    def __init__(self):
        """
        Generate authorization by public key and private key (https://tidbcloud.com/console/clusters)
        """
        self.digest_auth = _authorization()

    def create_manual_backup(self, project_id: str, cluster_id: str) -> dict:
        """
        Create manual backup
        :param project_id: The project id
        :param cluster_id: The dedicated cluster id
        :return: The backup id
        """
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters/{cluster_id}/backups"
        cur_date = datetime.datetime.now().strftime("%Y-%m-%d")
        data_for_backup = {"name": f"tidbcloud-backup-{cur_date}", "description": f"tidbcloud-backup-{cur_date}"}
        data_for_backup_json = json.dumps(data_for_backup)
        resp = requests.post(url=url,
                             data=data_for_backup_json,
                             auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}, Payload:{data_for_backup_json}")
        return _response(resp)

    def get_backup_info(self, project_id: str, cluster_id: str, backup_id: str) -> dict:
        """
        Get backup info by backup_id
        :param project_id: The project id
        :param cluster_id: The dedicated cluster id
        :param backup_id: Thr backup id
        :return: The backup detail
        """
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters/{cluster_id}/backups/{backup_id}"
        resp = requests.get(url=url,
                            auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        return _response(resp)

    def delete_backup(self, project_id: str, cluster_id: str, backup_id: str) -> dict:
        """
        Delete a backup for a cluster
        :param project_id: The project id
        :param cluster_id: The dedicated cluster id
        :param backup_id: The backup id
        :return:
        """
        url = f"{HOST}/api/v1beta/projects/{project_id}/clusters/{cluster_id}/backups/{backup_id}"
        resp = requests.delete(url=url,
                               auth=self.digest_auth)
        print(f"Method: {resp.request.method}, Request: {url}")
        return _response(resp)

    def create_restore_task(self, project_id: str, back_up_id: str, dedicated_config: dict) -> dict:
        """
        Create restore task
        :param project_id: The project id
        :param back_up_id: The backup id
        :param dedicated_config: The dedicated cluster config
        :return: The restore task id
        """
        try:
            tidb_size = dedicated_config["config"]["components"]["tidb"]["node_size"]
            tidb_quantity_range = dedicated_config["config"]["components"]["tidb"]["node_quantity"]
            tikv_size = dedicated_config["config"]["components"]["tikv"]["node_size"]
            tikv_quantity_range = dedicated_config["config"]["components"]["tikv"]["node_quantity"]
            tikv_storage_size_gib_range = dedicated_config["config"]["components"]["tikv"]["storage_size_gib"]
            tiflash_exist = False
            if dedicated_config["config"]["components"]["tiflash"] is not None:
                tiflash_exist = True
                tiflash_size = dedicated_config["config"]["components"]["tiflash"]["node_size"]
                tiflash_quantity_range = dedicated_config["config"]["components"]["tiflash"]["node_quantity"]
                tiflash_storage_size_gib_range = dedicated_config["config"]["components"]["tiflash"]["storage_size_gib"]
        except KeyError:
            print("cloud provider or region or available specifications not found!")
            raise
        url = f"{HOST}/api/v1beta/projects/{project_id}/restores"
        cur_date = datetime.datetime.now().strftime("%Y-%m-%d")
        if tiflash_exist:
            data_for_restore = \
                {
                    "backup_id": f"{back_up_id}",
                    "name": f"tidbcloud-restore-{cur_date}",
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
                                        },
                                    "tiflash":
                                        {
                                            "node_size": f"{tiflash_size}",
                                            "storage_size_gib": f"{tiflash_storage_size_gib_range}",
                                            "node_quantity": f"{tiflash_quantity_range}"
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
        else:
            data_for_restore = \
                {
                    "backup_id": f"{back_up_id}",
                    "name": f"tidbcloud-restore-{cur_date}",
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
        data_for_restore_json = json.dumps(data_for_restore)
        resp = requests.post(url=url, auth=self.digest_auth,
                             data=data_for_restore_json)
        print(f"Method: {resp.request.method}, Request: {url}, Payload:{data_for_restore_json}")
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


def usage_demo():
    print("-" * 88)
    print("Welcome to the TiDB Cloud API samples!")
    print("-" * 88)

    manage_backup = ManageBackup()
    dedicated_project_id = os.environ.get("DEDICATED_PROJECT_ID", None)
    dedicated_cluster_id = os.environ.get("DEDICATED_CLUSTER_ID", None)
    if dedicated_project_id is None or dedicated_cluster_id is None:
        print(
            "DEDICATED_PROJECT_ID or DEDICATED_CLUSTER_ID is None, you should set DEDICATED_PROJECT_ID and "
            "DEDICATED_CLUSTER_ID firstly.")
        raise Exception("DEDICATED_PROJECT_ID or DEDICATED_CLUSTER_ID not set!")

    print("1. Create a manual backup.")
    backup = manage_backup.create_manual_backup(dedicated_project_id, dedicated_cluster_id)
    try:
        sample_backup_id = backup["id"]
    except KeyError:
        print("backup id not found!")
        raise
    print()

    print(f"2. Show backup ( backup id : {sample_backup_id} ) process.")
    manage_backup.get_backup_info(dedicated_project_id, dedicated_cluster_id, sample_backup_id)
    print()

    # # tips: wait until backup completed before restore
    # sample_backup_id = "60346"
    # print("3. Restore backup data to a new cluster.")
    # dedicated_cluster_detail = manage_backup.get_cluster_by_id(dedicated_project_id, dedicated_cluster_id)
    # restore = manage_backup.create_restore_task(dedicated_project_id, sample_backup_id, dedicated_cluster_detail)
    # try:
    #     sample_restore_cluster_id = restore["cluster_id"]
    # except KeyError:
    #     print("the restore cluster id not found!")
    #     raise
    # print()
    #
    # print(f"4. Get the new cluster ( cluster id : {sample_restore_cluster_id} ) detail.")
    # manage_backup.get_cluster_by_id(dedicated_project_id, sample_restore_cluster_id)
    # print()

    # # tear down if necessary
    # print("If necessary , delete the backup.")
    # manage_backup.delete_backup(dedicated_project_id, dedicated_cluster_id, sample_backup_id)
    # print()
    #
    # print("If necessary , delete the init cluster.")
    # manage_backup.delete_cluster(dedicated_project_id, dedicated_cluster_id)

    print("Thanks for watching!")
    print("-" * 88)


if __name__ == "__main__":
    usage_demo()
