package auth

import "swis-api/users"

type AuthParams struct {
	BearerToken string `header:"X-Auth-Token" binding:"required"`
	User        users.User
}

var Params = AuthParams{
	BearerToken: "",
}
