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

var swives = []Frontend{
	{SiteName: "savla.dev", Title: "<h2>About Us</h2>", Description: "<p class=\"mb-3\">We are a group of open-minded Open Source enthusiasts.</p><p>We are interested in:<ul><li>IT Consulting</li><li>Software Development</li><li>IT Administration</li><li>Hosting</li></ul></p><p class=\"mb-3\">You can write to us at <a href=\"mailto: " + savlaDevMail + "\">" + savlaDevMail + "</a>.</p><p class=\"mb-3\">Visit us on <a href=\"" + savlaDevGithub + "\">GitHub</a>.</p>"},
}
