# Copyright(C) 2022 PingCAP. All Rights Reserved.
"""
Unit tests
"""

import pytest

from python.create_dedicated_cluster.main import CreateDedicatedCluster


class TestCreateCluster:
    def setup_method(self):
        self.create_cluster = CreateDedicatedCluster()
        self.project_id = None
        self.cluster_id = None

    @pytest.fixture()
    def setup_project(self):
        print("setup : get a project.")
        resp_body = self.create_cluster.get_all_projects()
        project_id = resp_body["items"][0]["id"]
        print(f"success. project id : {project_id}.")
        return project_id

    @pytest.fixture()
    def setup_region(self):
        print("setup : get a available region.")
        provider_region_specifications = self.create_cluster.get_provider_regions_specifications()
        items = provider_region_specifications["items"]
        for item in items:
            if item["cluster_type"] == "DEDICATED":
                return item
        return None

    @pytest.fixture()
    def setup_cluster(self, setup_region):
        print("setup : create a dedicated cluster.")
        resp_projects = self.create_cluster.get_all_projects()
        project_id = resp_projects["items"][0]["id"]
        resp_cluster = self.create_cluster.create_dedicated_cluster(project_id, setup_region)

        cluster_id = resp_cluster["id"]
        res = {
            "project_id": project_id,
            "cluster_id": cluster_id
        }
        self.project_id = project_id
        self.cluster_id = cluster_id
        print("success : create a dedicated cluster.")
        return res

    def teardown_method(self):
        print("tear down cluster. ")
        if self.project_id is not None and self.cluster_id is not None:
            self.create_cluster.delete_cluster(self.project_id, self.cluster_id)
        else:
            print("not necessary to tear down. ")

    def test_get_all_projects(self):
        print("test : get_all_projects.")
        resp_body = self.create_cluster.get_all_projects()

        # assert content
        assert resp_body["total"] >= 1

    def test_create_cluster(self, setup_project, setup_region):
        print("test : create cluster.")
        self.project_id = setup_project
        resp_body = self.create_cluster.create_dedicated_cluster(self.project_id, setup_region)

        # assert content
        assert "id" in resp_body

        self.cluster_id = resp_body["id"]

    def test_get_cluster_info(self, setup_cluster):
        print("test : get cluster info.")
        project_id = setup_cluster["project_id"]
        cluster_id = setup_cluster["cluster_id"]
        resp_body = self.create_cluster.get_cluster_by_id(project_id, cluster_id)

        # assert content
        assert resp_body["id"] == cluster_id
        assert resp_body["project_id"] == project_id


if __name__ == "__main__":
    pytest.main()
