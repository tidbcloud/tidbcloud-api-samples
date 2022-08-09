# Create a Dedicated Tier Cluster

This Go code sample shows you how to create a Dedicated Tier cluster on TiDB Cloud and get cluster's detail.

## Prerequisites

Before you start, make sure you have created the API key, configured environment variables, and installed the required Go version.

For details, refer to [Prerequisites for Go](../README.md#prerequisites).

## Run the sample

1. Clone the repository:

    ```
    git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
    ```

2. Change to the `create_dedicated_cluster` directory and run `make` to compile the code into an executable.

    ```bash
    cd go/create_dedicated_cluster/
    make default # GO11MODULE=on go build -o bin/create_cluster .
    ```

    After the compilation, the directory structure is similar to this:

    ```bash
    .
    ├── bin  # the binary is generated under this directory
    │   └── create_cluster
    ├── client.go  # http client
    ├── main.go  # entrypoint of this program
    ├── Makefile
    ├── README.md
    └── types.go  # all the structures we need for this demo
    ```

3. Execute the program.

    ```bash
    # export TIDBCLOUD_PUBLIC_KEY="<your public key>"
    # export TIDBCLOUD_PRIVATE_KEY="<your private key>"

    ./bin/create_cluster
    ```
