package users

// Low-level User struct with all user's details.
type User struct {
	// User ID as an unique identifier.
	ID string `json:"id" binding:"required" required:"true" readonly:"true"`

	// ID not used anymore as indexing is used differently now (searching by Name, index respects array implicit property) (legacy).
	Name string `json:"name" binding:"required" required:"true"`

	// Full Name of such user.
	FullName string `json:"full_name"`

	// User's given roles -- a role labels array.
	//Roles []roles.Role `json:"roles"`
	Roles []string `json:"roles"`

	// Access Control List. List of swapi modules to be accessed.
	ACL []string `json:"ACL"`

	// Presence/Absence boolean. If false, one is not allowed to log-in (token is rejected),
	// to interract with savla-dev infra in general (by default).
	Active bool `json:"active" default:false`

	// Unique token used for auth purposes, SHA512 preferred.
	TokenHash string `json:"token_hmac"`

	// GitHub account/profile name (used for SSH public keys importing).
	GitHubUser string `json:"github_username"`

	// Discord account/profile name.
	DiscordUser string `json:"discord_username"`

	// Spotify link to one's profile.
	SpotifyLink string `json:"spotify_link"`

	// Email address main, personal
	EmailMain string `json:"email_main" required:"true"`

	// Email alias in cloudflare email routing
	EmailAlias string `json:"email_alias"`

	// Country of origin -- to help maintain global contacts.
	Country string `json:"country"`

	// All Wireguard config objects -- an array.
	Wireguard []Wireguard `json:"wireguard_vpn"`

	// User's SSH public keys array.
	SSHKeys []string `json:"ssh_keys"`

	// User's GPG public keys array.
	GPGKeys []string `json:"gpg_keys"`

	// Important GDPR consent boolean -- if false, user's details should be omitted!
	// SEE more -- https://gdpr.eu/checklist/
	GDPRConsent bool `json:"gdpr_consent" default:false`
}

// Wireguard struct for the proper VPN connection purposes (prolly to be imported by vpn_gateway_server).
type Wireguard struct {
	// Unique device name (for such user).
	DeviceName string `json:"device_name"`

	// Wireguard public key.
	PublicKey string `json:"public_key"`

	// Wireguard private key.
	// TODO: should be encrypted?
	PrivateKey string `json:"private_key"`

	// User's private IP address.
	IPAddress string `json:"ip_address"`

	// Allowed IP address(es) list on the side of server (vpn_gateway_server).
	AllowedIPs []string `json:"allowed_ips"`

	// Is the user given permission to dial a connection?
	Permission bool `json:"permission" default:false`
}
