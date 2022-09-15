package swife

type Frontends struct {
	Frontends []Frontend `json:"frontends"`
}

type Frontend struct {
	// SiteName or hostname to get details for.
	SiteName string `json:"site_name"`

	// Site's title.
	Title string `json:"title"`

	// Site's description, possible on the frontpage.
	Description string `json:"description"`
}

const (
	// savla.dev info e-mail address
	savlaDevMail string = "info@savla.dev"

	// savla.dev github.com link
	savlaDevGithub string = "https://github.com/savla-dev"
)

// flush frontend
var swives = []Frontend{}
