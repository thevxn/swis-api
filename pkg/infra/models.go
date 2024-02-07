package infra

import (
	"time"
)

type Infrastructures struct {
	// Whole infrastructure object.
	Infrastructure Infrastructure `json:"infrastructure"`
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

	// Configuration system variables.
	Configuration Configuration `json:"configuration"`

	// Exported system facts from facter.
	Facts Facts `json:"facts"`

	// Provider is the name of the hosting company of such machine.
	Provider string `json:"provider"`

	// EpiresAt is a time of the expiration of such hosting service.
	ExpiresAt time.Time `json:"expires_at"`

	// Datacentre is the physical locality of such machine.
	Datacentre string `json:"datacentre"`
}

// facter-parsed system info, to-be-described further later TODO
type Facts struct {
	IsVirtual     bool   `json:"is_virtual"`
	KernelVersion string `json:"kernel_version"`

	MemoryTotalBytes int64 `json:"memory_total_bytes"`
	MemoryUsedBytes  int64 `json:"memory_used_bytes"`

	NetDomain         string `json:"net_domain"`
	NetHostname       string `json:"net_hostname"`
	NetFQDN           string `json:"net_fqdn"`
	NetPrimaryIP      string `json:"net_primary_ip"`
	NetPrimaryNetwork string `json:"net_primary_network"`

	OSArch    string `json:"os_arch"`
	OSFamily  string `json:"os_family"`
	OSSELinux bool   `json:"os_selinux_enabled"`

	ProcCount int `json:"proc_count"`

	SystemUptimeSec int64 `json:"system_uptime_sec"`

	Timestamp int64  `json:"timestamp"`
	Timezone  string `json:"timezone"`
}

// Configuration suits as a matrix for Ansible variables (as host_vars).
type Configuration struct {
	// ansible root vars
	AnsibleHost string `json:"ansible_host" yaml:"ansible_host"`
	AnsibleUser string `json:"ansible_user" yaml:"ansible_user"`
	Become      bool   `json:"become" yaml:"become" default:true`
	BecomeUser  string `json:"become_user" yaml:"become_user"`

	// base role
	// https://www.patorjk.com/software/taag/#p=display&f=ANSI%20Regular&t=stokrle
	BaseMotd        string `json:"base_motd" yaml:"ascii_art_motd" default:false`
	BaseDescription string `json:"base_description" yaml:"host_description"`

	// container role
	ContainerInstallk8sControl bool `json:"install_k8s_control_node" yaml:"install_k8s_control_node" default:false`

	// dialin-nas role
	DialInPresent   bool `json:"dialin_present" yaml:"dialin_present" default:false`
	AsteriskPresent bool `json:"asterisk_present" yaml:"asterisk_present" default:false`

	// dns role
	DNSServerPresent bool   `json:"dns_server_present" yaml:"dns_server_present" default:false`
	DMSServerType    string `json:"dns_server_type" yaml:"dns_server_type"`

	// ghar role
	RunnerPresent      bool   `json:"runner_present" yaml:"runner_present" default:false`
	RunnerAction       string `json:"runner_action" yaml:"runner_action"`
	RunnerVersion      string `json:"runner_version" yaml:"runner_version"`
	RunnerUser         string `json:"runner_user" yaml:"runner_user"`
	RunnerGroup        string `json:"runner_group" yaml:"runner_group"`
	RunnerConfigName   string `json:"runner_config_name" yaml:"runner_config_name"`
	RunnerConfigLabels string `json:"runner_config_labels" yaml:"runner_config_labels"`
	RunnerConfigToken  string `json:"runner_config_token" yaml:"runner_config_token"`

	// hyp vars
	IsHypervisor bool `json:"is_hypervisor" yaml:"is_hypervisor" default:false`

	// kpu role
	KPUPresent bool `json:"kpu_present" yaml:"kpu_present" default:false`

	// metrics role
	LokiPresent bool `json:"loki_present" yaml:"loki_present" default:false`
	LokiDockerTag string `json:"loki_image_tag" yaml:"loki_image_tag"`
	GrafanaPresent bool `json:"grafana_present" yaml:"grafana_present" default:false`
	GrafanaWebuiURL string `json:"grafana_webui_url" yaml:"grafana_webui_url"`
	GrafanaDockerVolume string `json:"grafana_docker_volume_name" yaml:"grafana_docker_volume_name"`
	GrafanaContainer string `json:"grafana_container_name" yaml:"grafana_container_name"`
	PrometheusPresent bool `json:"prometheus_present" yaml:"prometheus_present" default:false`
	PrometheusWebuiURL string `json:"prometheus_webui_url" yaml:"prometheus_webui_url"`
	PrometheusDockerVolume string `json:"prometheus_docker_volume_name" yaml:"prometheus_docker_volume_name"`
	PrometheusContainer string `json:"prometheus_container_name" yaml:"prometheus_container_name"`
	PrometheusDockerTag string `json:"prometheus_image_tag" yaml:"prometheus_image_tag"`
	PrometheusConfigDir string `json:"prometheus_config_dir" yaml:"prometheus_config_dir"`

	// net role
	NetWireguarded bool `json:"is_wireguarded" yaml:"is_wireguarded" default:false`

	// postfix role
	IsEdgeRelay bool `json:"is_edge_relay" yaml:"is_edge_relay" default:false`
	IsRelay bool `json:"is_relay" yaml:"is_relay" default:false`

	// proxy role
	IsBehindCf             bool   `json:"is_behind_cloudflare" yaml:"is_behind_cloudflare" default:false`
	NginxPresent           bool   `json:"nginx_present" yaml:"ngnix_present" default:false`
	NginxUseGeoIP          bool   `json:"use_geoip" yaml:"use_geoip" default:false`
	TraefikPresent         bool   `json:"traefik_present" yaml:"traefik_present" default:true`
	TraefikWebuiURL        string `json:"traefik_webui_url" yaml:"taefik_webui_url"`
	TraefikWebuiPort       int    `json:"traefik_webui_external_port" yaml:"traefik_webui_external_port"`
	TraefikDockerNet       string `json:"traefik_docker_network_name" yaml:"traefik_docker_network_name"`
	TraefikDockerTag       string `json:"traefik_docker_tag_version" yaml:"traefik_docker_tag_version"`
	TraefikDockarContainer string `json:"traefik_docker_container_name" yaml:"traefik_docker_container_name"`
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
