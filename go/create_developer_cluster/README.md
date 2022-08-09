# Create a Developer Tier Cluster

This Go code sample shows you how to create a Developer Tier cluster on TiDB Cloud and get cluster's detail.

## Prerequisites

Before you start, make sure you have created the API key, configured environment variables, and installed the required Go version.

For details, refer to [Prerequisites for Go](../README.md#prerequisites).

## Run the sample

1. Clone the repository:

    ```
    git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
    ```

2. Change to the `create_developer_cluster` directory and run `make` to compile the code into an executable.

    ```bash
    cd go/create_developer_cluster/
    make default # GO11MODULE=on go build -o bin/create_cluster .
    ```

    After the compilation, the directory structure is similar to this:

    ```bash
    .
    ├── bin  # the binary is generated under this directory
    │   └── create_cluster
    ├── clients.go  # http client
    ├── main.go  # entrypoint of this program
    ├── Makefile
    ├── README.md
    └── types.go  # all the structures we need for this demo
    ```

3. Execute the program.

    ```bash
    ./bin/create_cluster
    ```