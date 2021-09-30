package auth

const hasuraPrefix = "x-hasura-"
const organizationPrefix = hasuraPrefix + "org-"
const userPrefix = hasuraPrefix + "user-"

const defaultRoleHeaderName = hasuraPrefix + "role"

const orgIDHeaderName = organizationPrefix + "id"

// const orgNameHeaderName = organizationPrefix + "name"
const orgUsernameHeaderName = organizationPrefix + "username"

const userIDHeaderName = userPrefix + "id"

const roleIDHeaderName = hasuraPrefix + "role-id"
