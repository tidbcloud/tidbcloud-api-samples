# Copyright(C) 2022 PingCAP. All Rights Reserved.
"""
Unit tests
"""
import os

import pytest

from python.manage_backup.main import ManageBackup


class TestManageBackup:
    def setup_method(self):
        self.manage_backup = ManageBackup()
        self.project_id = os.environ.get("DEDICATED_PROJECT_ID", None)
        self.cluster_id = os.environ.get("DEDICATED_CLUSTER_ID", None)
        self.backup_id = None
        self.restore_cluster_id = None

    # # tips: tear down before backup is ready
    # def teardown_method(self):
    #     print("tear down backup. ")
    #     if self.project_id is not None and self.cluster_id is not None and self.backup_id is not None:
    #         self.manage_backup.delete_backup(self.project_id, self.cluster_id, self.backup_id)
    #     else:
    #         print(" not necessary to tear down. ")
    #     print("tear down restore cluster. ")
    #     if self.project_id is not None and self.restore_cluster_id is not None:
    #         self.manage_backup.delete_cluster(self.project_id, self.restore_cluster_id)
    #     else:
    #         print(" not necessary to tear down. ")

    @pytest.fixture()
    def setup_backup(self):
        print("setup : create a backup.")
        resp_body = self.manage_backup.create_manual_backup(self.project_id, self.cluster_id)
        self.backup_id = resp_body["id"]
        return resp_body["id"]

    def test_create_manual_backup(self):
        print("test : create manual backup.")
        resp_body = self.manage_backup.create_manual_backup(self.project_id, self.cluster_id)

        # assert content
        assert "id" in resp_body

        self.backup_id = resp_body["id"]

    # # tips: test before backup is ready
    # def test_create_restore_task(self, setup_backup):
    #     print("test : create restore task.")
    #     resp_body = self.manage_backup.create_restore_task(self.project_id, self.backup_id)
    #
    #     # assert content
    #     assert "cluster_id" in resp_body
    #
    #     self.restore_cluster_id = resp_body["cluster_id"]


if __name__ == "__main__":
    pytest.main()
