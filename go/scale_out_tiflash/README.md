# Scale Out a TiFlash Node

This Go code sample shows you how to scale out a TiFlash node for your cluster and view the scale-out progress.

## Prerequisites

Before you start, make sure you have created the API key, configured environment variables, and installed the required Go version.

For details, refer to [Prerequisites for Go](../README.md#prerequisites).

## Run the sample

1. Clone the repository:

    ```
    git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
    ```

2. Change to the `scale_out_tiflash` directory and run `make` to compile the code into an executable.

    ```bash
    cd go/scale_out_tiflash/
    make default # GO11MODULE=on go build -o bin/scale_out_tiflash .
    ```

    After the compilation, the directory structure is similar to this:

    ```bash
    .
    ├── bin  # the binary is generated under this directory
    │   └── scale_out_tiflash
    ├── client.go  # http client
    ├── main.go  # entrypoint of this program
    ├── Makefile
    ├── README.md
    └── types.go  # all the structures we need for this demo
    ```

3. Execute the program.

    ```bash
    ./bin/scale_out_tiflash
    ```
