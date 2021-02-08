package main

import (
	"encoding/binary"
	"errors"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
)

// this function doesn't handle all the edge cases, but unless we have > 256 servers it should be fine
func AllocateIP(subnet net.IPNet, taken []net.IPAddr) net.IPAddr {
	ip := make(net.IP, 4)

	if len(taken) == 0 {
		binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(subnet.IP.To4())+1)
		return net.IPAddr{IP: ip}
	} else {
		maxIP := binary.BigEndian.Uint32(taken[0].IP.To4())
		for _, ip := range taken {
			ipInt := binary.BigEndian.Uint32(ip.IP.To4())
			if ipInt > maxIP {
				maxIP = ipInt
			}
		}
		binary.BigEndian.PutUint32(ip, maxIP+1)
		return net.IPAddr{IP: ip}
	}
}

func Save() error {
	out, _ := yaml.Marshal(servers)
	return ioutil.WriteFile(os.Args[2], out, 0644)
}

func AddServer(serverType string, metadata map[string]string) (Server, error) {
	var ipRange net.IPNet
	var found bool
	for _, configType := range config.Types {
		if configType.Name == serverType {
			ipRange = *configType.Cidr.IPNet
			found = true
			break
		}
	}

	if !found {
		return Server{}, errors.New("server type not found")
	}

	takenIPs := make([]net.IPAddr, 0)
	for _, server := range servers {
		if server.Type == serverType {
			takenIPs = append(takenIPs, net.IPAddr{IP: server.IP.IP})
		}
	}

	ip := AllocateIP(ipRange, takenIPs)

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return Server{}, err
	}

	server := Server{
		IP: IPAddr{&ip},
		AllowedIPs: []IPNet{
			{
				IPNet: &net.IPNet{IP: config.HubIP.IP,
					Mask: []byte{255, 255, 255, 255},
				},
			},
			{
				IPNet: &ipRange,
			},
		},
		Type:       serverType,
		Metadata:   metadata,
		PublicKey:  key.PublicKey().String(),
		PrivateKey: key.String(),
	}

	servers = append(servers, server)

	err = Save()

	if err != nil {
		return Server{}, err
	}

	err = wgClient.ConfigureDevice(config.Interface, wgtypes.Config{
		Peers: []wgtypes.PeerConfig{
			{
				PublicKey: key.PublicKey(),
				AllowedIPs: []net.IPNet{
					{
						IP:   ip.IP,
						Mask: []byte{255, 255, 255, 255}, // = /32
					},
				},
			},
		},
	})

	if err != nil {
		return Server{}, err
	}

	return server, nil
}

func DeleteServer(ip IPAddr) error {
	index := -1
	for i, server := range servers {
		if server.IP.String() == ip.String() {
			index = i
		}
	}
	if index == -1 {
		return errors.New("server not found")
	}

	servers = append(servers[:index], servers[index+1:]...)
	newPeers := make([]wgtypes.PeerConfig, len(servers))

	for i := 0; i < len(servers); i++ {
		key, _ := wgtypes.ParseKey(servers[i].PrivateKey)

		newPeers[i] = wgtypes.PeerConfig{
			PublicKey: key.PublicKey(),
			AllowedIPs: []net.IPNet{
				{
					IP:   servers[i].IP.IP,
					Mask: []byte{255, 255, 255, 255}, // = /32
				},
			},
		}
	}

	err := wgClient.ConfigureDevice(config.Interface, wgtypes.Config{
		Peers: newPeers,
	})

	if err != nil {
		return err
	}

	return Save()
}
