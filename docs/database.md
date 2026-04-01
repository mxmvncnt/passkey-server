# Database
General documentation relating to the database.Database

## Schema
All that is required to run this project is a PostgreSQL instance with the table defined in [schema.sql](../database/schema.sql)

### Creating the new table
You will need to create a new SQL table in your Postgres schema to use this project. Before you create the table defined
in the `database/schema.sql` file, you need to first change the values that will relate to your Users table and user ID
columns. For example:

```sql
CREATE TABLE webauthn_credentials
(
    id                   BYTEA PRIMARY KEY,
    user_id              UUID        NOT NULL REFERENCES users_table (users_table_id_column) ON DELETE CASCADE, /* <--- CHANGE THIS!!!! */
    public_key           BYTEA       NOT NULL,
    /*  rest of schema ... */
);

CREATE INDEX idx_webauthn_credentials_user_id ON webauthn_credentials (user_id);
```

In the example above, I replaced the `users` and `id` values with `users_table` and `users_table_id_column`, where of
course you will need to replace this with your own values.

## Connecting to your existing database
To connect the database to your existing project, your Users table **must** use a `UUID` as the identity column (index).
If you do not use `UUID`s as an identity column, you can always add a column to your existing user table.