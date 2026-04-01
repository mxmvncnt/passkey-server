# Passkey Server
This project aims to make passkey rollout simpler by providing a plug and play solution to add passkeys into applications.

This project is made to integrate seamlessly in your existing Postgres database to make it easier to install.

## Dev Setup
This project uses [SQLC](https://docs.sqlc.dev/en/stable/overview/install.html) for generating the schema definitions to
 be used in the code.
Go 1.26 and above is required.

After everything is installed, you can proceed.
1. Install dependencies
    ```bash
    go mod tidy
    ```
2. (If updating database schemas or queries) Re-generate SQLC code
    ```bash
    sqlc generate
    ```
3. Create a config
    a. Copy [config/example.env](./config/example.env)
    b. Rename it to something you want (for example, `dev.env`)
    c. Edit the values to match your preferences. For more details, look at the [config docs](docs/config.md)
3. Run the project
    ```
    go run .
    ```

## Docs
Documentation is located in the [documentation directory](docs/index.md).