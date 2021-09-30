# hasura-auth-webhook

Utility server to be used with Hasura projects. It stores the users on its own database. It communicates with Hasura through the webhook setting and returning only the user ID and the user role.

## Usage

Hasura instance must point to this server for authentication. This can be done by setting the following environments keys on your hasura instance.

```
HASURA_GRAPHQL_AUTH_HOOK=<this-server>
HASURA_GRAPHQL_AUTH_HOOK_MODE=POST
```

You can start the server with

```bash
task start
```
