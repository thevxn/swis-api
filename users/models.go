package users

type Users struct {
	Users []User `json:"users"`
}

type User struct {
	// ID not used anymore as indexing is used differently now (searching by Name, index respects array implicit property)
	//ID       	string 		`json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Roles       []string    `json:"roles"`
	TokenBase64 string      `json:"token_base64"`
	GitHubUser  string      `json:"github_username"`
	Wireguard   []Wireguard `json:"wireguard_vpn"`
	SSHKeys     []string    `json:"ssh_keys"`
	GPGKeys     []string    `json:"gpg_keys"`
	// SEE more -- https://gdpr.eu/checklist/
	GDPRConsent bool `json:"gdpr_consent" default:false`
}

type Wireguard struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	// user's private IP address
	IPAddress string `json:"ip_address"`
	// IP address list on the side of server (frank)
	AllowedIPs []string `json:"allowed_ips"`
	Permission bool     `json:"permission" default:false`
	DeviceName string   `json:"device_name"`
}

// flush users at start -- see Makefile, import_prod target, and .data/users
var users = []User{} //equivalent to Users.Users{}
