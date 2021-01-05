package main

import (
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

var config Config
var wgClient *wgctrl.Client
var device *wgtypes.Device
var servers []Server

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("usage: ./central [config file] [servers file]")
	}

	configData, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalln("failed to read config file", err)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatalln("failed to parse config file", err)
	}

	serverData, err := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatalln("failed to read server list file", err)
	}

	err = yaml.Unmarshal(serverData, &servers)
	if err != nil {
		log.Fatalln("failed to parse server list file", err)
	}

	wgClient, err = wgctrl.New()
	if err != nil {
		log.Fatalln("failed to open wireguard interface", err)
	}

	device, err = wgClient.Device(config.Interface)
	if err != nil {
		log.Fatalln("could not find interface "+config.Interface, err)
	}
	log.Printf("opened wg interface %s (pubkey = %s)", config.Interface, device.PublicKey)

	for _, serverType := range config.Types {
		log.Printf("found type %s, assigned range %s", serverType.Name, serverType.Cidr)
	}

	for _, server := range servers {
		log.Printf("found server %s of type %s", server.IP, server.Type)
	}

	StartServer()
}
