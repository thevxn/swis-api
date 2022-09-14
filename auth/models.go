package auth

import (
	//"swis-api/roles"
	"swis-api/users"
)

type AuthParams struct {
	// Token string to load and verify against internal 'users' module, authentication.
	BearerToken string `header:"X-Auth-Token" binding:"required" verification:"required"`

	// User object to add to server context.
	User users.User

	// Roles string array/objects for ACL authorization.
	//Roles       roles.Roles
	Roles []string
}

var Params = AuthParams{
	// Wipe Token string at every request not to allow token forgery.
	BearerToken: "",
}
