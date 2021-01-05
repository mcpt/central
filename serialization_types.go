package main

import (
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net"
)

type IPAddr struct {
	*net.IPAddr
}

func (ip *IPAddr) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	parsed := net.ParseIP(s)
	if parsed == nil {
		return errors.New("invalid ip")
	}

	*ip = IPAddr{&net.IPAddr{IP: parsed}}
	return nil
}


func (ip IPAddr) MarshalYAML() (interface{}, error) {
	return ip.String(), nil
}

func (ip IPAddr) MarshalJSON() ([]byte, error) {
	return json.Marshal(ip.String())
}

type IPNet struct {
	*net.IPNet
}

func (ip *IPNet) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	_, parsed, err := net.ParseCIDR(s)
	if err != nil {
		return err
	}

	*ip = IPNet{parsed}
	return nil
}

func (ip IPNet) MarshalYAML() (interface{}, error) {
	return ip.String(), nil
}

func (ip IPNet) MarshalJSON() ([]byte, error) {
	return json.Marshal(ip.String())
}

type Secret struct {
	Secret []byte
}

func (secret *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	err := unmarshal(&s)
	if err != nil {
		return err
	}

	decoded, err := hex.DecodeString(s)
	if err != nil {
		return err
	}

	*secret = Secret{decoded}
	return nil
}

func (secret Secret) Compare(other Secret) bool {
	return subtle.ConstantTimeCompare(other.Secret, secret.Secret) == 1
}