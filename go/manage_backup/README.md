# Manage Backup and Restoration

This Go code sample shows you how to mange backup and restoration, including the following tasks:

- Create a manual backup
- Get the backup info
- Restore backup data to a new cluster
- Delete a backup

## Prerequisites

Before you start, make sure you have created the API key, configured environment variables, and installed the required Go version.

For details, refer to [Prerequisites for Go](../README.md#prerequisites).

Note that you can only create backups for a Dedicated Tier cluster.

### Run the sample

1. Clone the repository:

    ```
    git clone https://github.com/tidbcloud/tidbcloud-api-samples.git
    ```

2. Change to the `manage_backup` directory, run `make` to compile the code into an executable.

    ```bash
    cd go/manage_backup/
    make default # GO11MODULE=on go build -o bin/manage_backup .
    ```

    After the compilation, the directory structure is similar to this:

    ```bash
    .
    ├── bin  # the binary is generated under this directory
    │   └── manage_backup
    ├── client.go  # http client
    ├── main.go  # entrypoint of this program
    ├── Makefile
    ├── README.md
    └── types.go  # all the structures we need for this demo
    ```

3. Execute the program.

    ```bash
    ./bin/manage_backup
    ```
