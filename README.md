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

## Env Variables

```
DATABASE_URL

API_PUBLIC_HOSTNAME
API_PUBLIC_PORT
API_INTERNAL_HOSTNAME
API_INTERNAL_PORT

ADMIN_USER_EMAIL
ADMIN_USER_PASSWORD

PROVIDERS_EMAIL_ENABLED
PROVIDERS_EMAIL_JWT_ACCESS_SECRET
PROVIDERS_EMAIL_JWT_REFRESH_SECRET
PROVIDERS_EMAIL_WEBHOOKS_RECOVERY_PASSWORD_EVENT
PROVIDERS_EMAIL_WEBHOOKS_REGISTER_EVENT
PROVIDERS_MAGICLINK_ENABLED
PROVIDERS_MAGICLINK_WEBHOOKS_LOGIN_EVENT
```
