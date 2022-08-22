package swife

type Frontend struct {
	// SiteName or hostname to get details for.
	SiteName string `json:"site_name"`

	// Site's title.
	Title string `json:"title"`

	// Site's description, possible on the frontpage.
	Description string `json:"description"`
}

const (
	savlaDevMail   string = "info@savla.dev"
	savlaDevGithub string = "https://github.com/savla-dev"
)

// flush frontend
var swives = []Frontend{}
