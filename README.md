# Passkey Server
This project aims to make passkey rollout simpler by providing a plug and play solution to add passkeys into applications.

This project is made to integrate seamlessly in your existing Postgres database to make it easier to install.

## Warning
This project is intended to be used as a micro-service, which means it should not be exposeddirectly to the users, at
least for the time being.

### Woke Software Disclaimer
This is woke software because I think gay people should exist.
As such, this is an application to be added on [Lunduke's SoftwarePoliticsTracker](https://github.com/BryanLunduke/SoftwarePoliticsTracker)

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
    1. Copy [config/example.env](./config/example.env)
    2. Rename it to something you want (for example, `dev.env`)
    3. Edit the values to match your preferences. For more details, look at the [config docs](docs/config.md)
3. Run the project
    ```
    (set -a source config/example.env set +a go run .)
    ```

## Docs
Documentation is located in the [documentation directory](docs/index.md).