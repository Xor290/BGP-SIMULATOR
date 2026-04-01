package models

type PeerVars struct {
	PeerIP      string `yaml:"peer_ip"`
	RemoteASN   uint32 `yaml:"remote_asn"`
	Description string `yaml:"description,omitempty"`
	Password    string `yaml:"password,omitempty"`
	Enabled     bool   `yaml:"enabled"`
}

type PrefixVars struct {
	Prefix    string `yaml:"prefix"`
	NextHop   string `yaml:"next_hop"`
	LocalPref uint   `yaml:"local_pref"`
}

type PlaybookVars struct {
	LocalASN uint32       `yaml:"local_asn"`
	RouterID string       `yaml:"router_id"`
	Hostname string       `yaml:"hostname"`
	Peers    []PeerVars   `yaml:"peers"`
	Prefixes []PrefixVars `yaml:"prefixes"`
}
