package infra

type Infrastructures struct {
	Infrastructure Infrastructure `json:"infrastructure"`
}

type Infrastructure struct {
	//Users		[]User
	//Groups	[]Group
	Hosts    []Host    `json:"hosts"`
	Networks []Network `json:"networks"`
}

type Network struct {
	Hash      string `json:"hash"`
	Name      string `json:"network_name"`
	Address   string `json:"network_address"`
	Interface string `json:"interface"`
	CIDRBlock string `json:"network_cidr_block"`
}

/*
 * IP address alocation (CIDR) count + possible supernetwork "expandability"
 * https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing#IPv4_CIDR_blocks
 *
 * /24 -- 256-2
 * /25 -- 128-2
 * /26 -- 64-2 (bigger lan networks)
 * /27 -- 32-2 (small lan networks)
 *
 */

/*
var networks = []Network{
	{Name: "squabbit virbr0 network", Address: "10.4.5.0", CIDRBlock: "27"},
	{Name: "-vacant virbr32 network", Address: "10.4.5.32", CIDRBlock: "27"},
	{Name: "-vacant virbr32 network", Address: "10.4.5.64", CIDRBlock: "27"},
	{Name: "-vacant network", Address: "10.4.5.96", CIDRBlock: "27"},
	{Name: "-vacant network", Address: "10.4.5.128", CIDRBlock: "27"},
	{Name: "VPN external dish network", Address: "10.4.5.160", CIDRBlock: "27"},
	{Name: "VPN private intranet", Address: "10.4.5.192", CIDRBlock: "27"},
	{Name: "VPN client route-all-traffic network", Address: "10.4.5.224", CIDRBlock: "27"},
}
*/

// High-level struct for batch []Host array importing.
type Hosts struct {
	Hosts []Host `json:"hosts"`
}

// Low-level Host model for a generic machine.
type Host struct {
	// Unique hash/ID to link to such host record.
	ID string `json:"id" binding:"required"`

	// Node hostname without its domain.
	HostnameShort string `json:"hostname_short" binding:"required"`

	// Node hostname as FQDN-formatted.
	HostnameFQDN string `json:"hostname_fqdn" binding:"required"`

	// Brief node's description -- shown in node's MOTD on remote login.
	Description string `json:"description"`

	// Host's default domain name (e.g. savla.su as internal domain name).
	Domain string `json:"domain" `

	// Ansible roles to be applied to such host.
	Roles []string `json:"roles"`

	// Important network-related IP addresses to be assigned to such host (e.g. public interface address, wireguard interface address etc).
	IPAddress []string `json:"ip_address"`

	// Children of such machine -- should use machines' hashes for proper linking.
	Child []string `json:"children"`
}

var infrastructure = Infrastructure{}

//var hosts = []Host{}
//var networks = []Network{}

// Hyper struct to model hypervisor machine
type Hyper struct {
	Host
}

// Virtual struct to model virtual machine
type Virtual struct {
	Host
}

var squabbitVMs = []Host{
	{HostnameShort: "stokrle", Domain: "savla.su", Roles: []string{"build", "deploy"}, IPAddress: []string{"10.4.5.55/25"}},
}

// demo Hosts data
/*
var hosts = Hosts{
	{Hostname: "squabbit", Domain: "savla.su", Role: "hypervisor", IPAddress: []string{"10.4.5.1/25", "10.4.5.129/25"}, VMs: []Host},
}
*/
