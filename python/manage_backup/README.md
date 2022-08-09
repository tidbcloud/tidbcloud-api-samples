# Manage Backup and Restoration

This Python code sample shows you how to mange backup and restoration, including the following tasks:

- Create a manual backup
- Get the backup info
- Restore backup data to a new cluster (commented out by default)
- Delete a backup (commented out by default)

## Prerequisites

Before you start, make sure you have created the API key, get the Dedicated Tier cluster information, configured environment variables, and installed the required Python version.

For details, refer to [Prerequisites for Python](../README.md#prerequisites).

Note that you can only create backups for a Dedicated Tier cluster.

## Run the sample

To run the complete code sample, take the following steps:

1. Clone the repository:

   ```
   git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
   ```

2. Install the dependencies using pip:

   ```shell
   cd tidbcloud-api-sample-test/python/manage_backup
   pip install -r requirements.txt # Might be "pip3" depending on your Python installation
   ```

3. Run the sample:

    ```shell
    # export TIDBCLOUD_PUBLIC_KEY="<your public key>"
    # export TIDBCLOUD_PRIVATE_KEY="<your private key>"
    # export DEDICATED_PROJECT_ID="<your dedicated project id>"
    # export DEDICATED_CLUSTER_ID="<your dedicated cluster id>"

    python main.py # Might be "python3" depending on your Python installation
    ```
