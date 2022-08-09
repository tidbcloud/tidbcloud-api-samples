<!-- Copyright(C) 2022 PingCAP. All Rights Reserved. -->

# TiDB Cloud API Samples Test

__*Samples are currently in preview.*__

This repository stores the code samples for TiDB Cloud API. You can find samples that use the TiDB Cloud API to interact with TiDB Cloud resources.

Currently, the code samples are available in the following programming languages:

- [Go](./go)
- [Python](./python)

With these samples, you can use API to perform the following tasks:

- Create a Dedicated Tier cluster ([Go](./go/create_dedicated_cluster) | [Python](./python/create_dedicated_cluster))
- Create a Developer Tier cluster ([Go](./go/create_developer_cluster) | [Python](./python/create_developer_cluster))
- Manage backups ([Go](./go/manage_backup) | [Python](./python/manage_backup))
- Scale out a TiFlash node ([Go](./go/scale_out_tiflash) | [Python](./python/scale_out_tiflash))

To learn more details on these code samples, see [the code samples walk-through](https://docs.pingcap.com/tidbcloud/api/v1beta#section/Get-Started/Code-samples) in TiDB Cloud documentation. For more information on getting started with the TiDB Cloud API, see [TiDB Cloud API Documentation](https://docs.pingcap.com/tidbcloud/api/v1beta).

## Quick start

### Create a TiDB Cloud account

To use the samples in this repository, you must have a TiDB Cloud account. If you do not have an account yet, click [here](https://tidbcloud.com/signup) to create one.

### Get an API key

To access the API, you'll need an API key. To create an API key, log in to your [TiDB Cloud console](https://tidbcloud.com/console). Navigate to the **Organization Settings** page, and create an API key.

An API key contains a public key and a private key. Copy and save the private key in a secure location. After leaving this page, you will not be able to get the full private key again.

For more details on API key, see [API Key Management](https://docs.pingcap.com/tidbcloud/api/v1beta#section/Authentication/API-Key-Management).

### (Optional) Get project ID and cluster ID

To manage an existing Dedicated Tier cluster, such as creating a backup and scaling out a node, you need to obtain `dedicated_project_id` and `dedicated_cluster_id`.

To obtain `dedicated_project_id`, you can call the [List all accessible projects](https://docs.pingcap.com/tidbcloud/api/v1beta#tag/Project/operation/ListProjects) endpoint. The `items.id` in response is the `dedicated_project_id`.

```shell
curl --digest \
    --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \
    --request GET \
    --url 'https://api.tidbcloud.com/api/v1beta/projects?page=1&page_size=10'
```

To obtain `dedicated_cluster_id`, you can call the [List all clusters in a project](https://docs.pingcap.com/tidbcloud/api/v1beta#tag/Cluster/operation/ListClustersOfProject) endpoint and replace `{dedicated_project_id}` with your project ID obtained from the preceding command. The `items.id` in response is the `dedicated_cluster_id`.

```shell
curl --digest \
    --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \
    --request GET \
    --url 'https://api.tidbcloud.com/api/v1beta/projects/{dedicated_project_id}/clusters?page=1&page_size=10'
```

### Configure environment variables

Store your API key in environment variables. These environment variables will be used by all code samples in this repository.

For example, in UNIX-like OS, you can do the following:

```shell
export TIDBCLOUD_PUBLIC_KEY="<your public key>"
export TIDBCLOUD_PRIVATE_KEY="<your private key>"
```

Store your Dedicated Tier information in environment variables. These environment variables will be used by `manage_backup` and `scale_out_tiflash` code samples in this repository.

```shell
export DEDICATED_PROJECT_ID="<your dedicated project id>"
export DEDICATED_CLUSTER_ID="<your dedicated cluster id>"
```

### Next

Select a programming language of your choice and run the code samples:

- [Go](./go)
- [Python](./python)

## ⚠️ Important

Running these code examples might result in charges to your TiDB Cloud account. For more information, see [TiDB Cloud Billing](https://docs.pingcap.com/tidbcloud/tidb-cloud-billing).

Some samples might create resources that have long-term costs with services. Some examples might modify or delete resources. It is your responsibility to do the following:

* Be aware of the resources that these examples create or delete.
* Be aware of the costs that might be charged to your account as a result.
* Back up your important data.
