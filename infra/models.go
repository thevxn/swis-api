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

type Hosts struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	Hash      string   `json:"hash"`
	Hostname  string   `json:"hostname"`
	Domain    string   `json:"domain"`
	Roles     []string `json:"roles"`
	IPAddress []string `json:"ip_address"`
	Facter    []string `json:"facter"`
	VMs       []string `json:"virtual_machines"`
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
	{Hostname: "stokrle", Domain: "savla.su", Roles: []string{"build", "deploy"}, IPAddress: []string{"10.4.5.55/25"}},
}

// demo Hosts data
/*
var hosts = Hosts{
	{Hostname: "squabbit", Domain: "savla.su", Role: "hypervisor", IPAddress: []string{"10.4.5.1/25", "10.4.5.129/25"}, VMs: []Host},
}
*/
