# Copyright(C) 2022 PingCAP. All Rights Reserved.
"""
Unit tests
"""
import os

import pytest

from python.scale_out_tiflash.main import ScaleOutTiFlash


class TestScaleOutTiFlash:
    def setup_method(self):
        self.scale_out_tiflash = ScaleOutTiFlash()
        self.project_id = os.environ.get("DEDICATED_PROJECT_ID", None)
        self.cluster_id = os.environ.get("DEDICATED_CLUSTER_ID", None)

    def test_modify_cluster(self):
        print("test : modify cluster.")
        dedicated_config = self.scale_out_tiflash.get_cluster_by_id(self.project_id, self.cluster_id)
        resp_body = self.scale_out_tiflash.modify_cluster(self.project_id, self.cluster_id, dedicated_config)

        assert resp_body == {}


if __name__ == "__main__":
    pytest.main()
