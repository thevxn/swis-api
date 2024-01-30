package infra

type Infrastructures struct {
	// Whole infrastructure object.
	Infrastructure []Infrastructure `json:"infrastructure"`
}

type Infrastructure struct {
	// Domains linked to such infrastructure.
	Domains []Domain `json:"domains"`

	// Hosts/nodes linked to such infrastructure.
	Hosts []Host `json:"hosts"`

	// Networks linked to such infrastructure.
	Networks []Network `json:"networks"`
}

type Domain struct {
	// Unique domain identifier.
	ID string `json:"domain_id"`

	// Fully qualified domain name.
	FQDN string `json:"domain_fqdn"`

	// Domain's owner (user)name.
	Owner string `json:"domain_owner"`

	// Domain's expiration date
	ExpirationDate string `json:"expiration_date"`

	// Name of the current registrar
	RegistrarName string `json:"registrar_name"`

	// Private details (e.g. addresses, phone number etc in WHOIS)
	PrivateDetails bool `json:"private_details" default:false`

	// Cloudflare Zone ID
	CfZoneID string `json:"cf_zone_id"`
}

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
	Domain string `json:"domain"`

	// Ansible roles to be applied to such host.
	Roles []string `json:"roles"`

	// Important network-related IP addresses to be assigned to such host (e.g. public interface address, wireguard interface address etc).
	IPAddress []string `json:"ip_address"`

	// Children of such machine -- should use machines' hashes for proper linking.
	Children []string `json:"children"`

	// Wireguarded bool indicates that the host is part of the core network.
	Wireguarded bool `json:"wireguarded"`

	// Exported system facts from facter.
	Facts Facts `json:"facts"`

	// Provider is the name of the hosting company of such machine.
	Provider string `json"provider"`

	// EpiresAt is a time of the expiration of such hosting service.
	ExpiresAt time.Time `json:"expires_at"`

	// Datacentre is the physical locality of such machine.
	Datacentre string `json:"datacentre"`
}

// facter-parsed system info, to-be-described further later TODO
type Facts struct {
	IsVirtual     bool   `json:"is_virtual"`
	KernelVersion string `json:"kernel_version"`

	MemoryTotalBytes int `json:"memory_total_bytes"`
	MemoryUsedBytes  int `json:"memory_used_bytes"`

	NetDomain         string `json:"net_domain"`
	NetHostname       string `json:"net_hostname"`
	NetFQDN           string `json:"net_fqdn"`
	NetPrimaryIP      string `json:"net_primary_ip"`
	NetPrimaryNetwork string `json:"net_primary_network"`

	OSArch    string `json:"os_arch"`
	OSFamily  string `json:"os_family"`
	OSSELinux bool   `json:"os_selinux_enabled"`

	ProcCores int `json:"proc_cores"`

	SystemUptimeSec int `json:"system_uptime_sec"`

	Timestamp int64  `json:"timestamp"`
	Timezone  string `json:"timezone"`
}

type Network struct {
	// Unique network's identifier
	Hash string `json:"hash"`

	// Network name, verbose ID.
	Name string `json:"network_name"`

	// Network IP address.
	Address string `json:"network_address"`

	// Interface(s) of such network.
	Interface string `json:"interface"`

	// CIDR block of netmask.
	CIDRBlock string `json:"network_cidr_block"`
}
