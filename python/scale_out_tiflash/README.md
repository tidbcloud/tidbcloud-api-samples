# Scale Out a TiFlash Node

This Python code sample shows you how to scale out a TiFlash node for your cluster and view the scale-out progress.

## Prerequisites

Before you start, make sure you have created the API key, get the Dedicated Tier cluster information, configured environment variables, and installed the required Python version.

For details, refer to [Prerequisites for Python](../README.md#prerequisites).

## Run the sample

To run the complete code sample, take the following steps.

1. Clone the repository:

    ```
    git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
    ```

2. Install the dependencies using pip:

    ```shell
    cd tidbcloud-api-sample-test/python/scale_out_tiflash
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
