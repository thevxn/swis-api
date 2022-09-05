package users

import "swis-api/roles"

// High-level Users struct mainly for users restore batch importing.
type Users struct {
	// An array of User objects
	Users []User `json:"users"`
}

// Low-level User struct with all user's details.
type User struct {
	// ID not used anymore as indexing is used differently now (searching by Name, index respects array implicit property).
	//ID       	string 		`json:"id"`
	Name string `json:"name" binding:"required"`

	// Full Name of such user.
	FullName string `json:"full_name"`

	// User's given roles -- a role labels array.
	Roles []roles.Role `json:"roles"`

	// Presence/Absence boolean. If 'absent', one is not allowed to log in, to interract with savla-dev infra in general (by default).
	State string `json:"state" default:"absent"`

	// Unique token used for auth purposes, base64'd.
	TokenBase64 string `json:"token_base64"`

	// GitHub account/profile name (used for SSH public keys importing).
	GitHubUser string `json:"github_username"`

	// Discord account/profile name.
	DiscordUser string `json:"discord_username"`

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

// Flush users pointer (sic!) at start -- see Makefile, import_prod[...] target, and .data/users.
var users = []User{} //equivalent to Users.Users{}
