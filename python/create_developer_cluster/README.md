# Create a Developer Tier Cluster

This Python code sample shows you how to create a Developer Tier cluster on TiDB Cloud and get cluster's detail.

## Prerequisites

Before you start, make sure you have created the API key, configured environment variables, and installed the required Python version.

For details, refer to [Prerequisites for Python](../README.md#prerequisites).

## Run the sample

To run the complete code sample, take the following steps.

1. Clone the repository:

   ```
   git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
   ```

2. Install the dependencies using pip:

   ```shell
   cd tidbcloud-api-sample-test/python/create_developer_cluster
   pip install -r requirements.txt # Might be "pip3" depending on your Python installation
   ```

3. Run the sample:

    ```shell
    # export TIDBCLOUD_PUBLIC_KEY="<your public key>"
    # export TIDBCLOUD_PRIVATE_KEY="<your private key>"

    python main.py # Might be "python3" depending on your Python installation
    ```
