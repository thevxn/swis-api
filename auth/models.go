package auth

type AuthParams struct {
	BearerToken string `header:"X-Auth-Bearer"`
}

var (
	Params = AuthParams{
		BearerToken: "",
	}
)
