package main

type Type struct {
	Name string
	Cidr IPNet
}

type Config struct {
	HubIP     IPAddr
	Interface string
	Listen    string
	Types     []Type
	Secret    Secret
}

type Server struct {
	IP         IPAddr
	Type       string
	AllowedIPs []IPNet
	Metadata   map[string]string
	PublicKey  string
	PrivateKey string
}
