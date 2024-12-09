package api

import "strings"

// NetworkAddressSet represents an address set for OVN.
// Refer to doc/network-acls.md for details.
//
// swagger:model
//
// API extension: network_address_set.
type NetworkAddressSet struct {
	// Name of the address set
	// Example: "core_services"
	Name string `json:"name" yaml:"name"`

	// Addresses included in the set (IPv4/IPv6/CIDR)
	// Example: ["10.0.0.5", "192.168.0.1/24"]
	Addresses []string `json:"addresses" yaml:"addresses"`

	// Mapping of key-value pairs for custom use
	// Example: {"prod": "false"}
	ExternalIDs map[string]string `json:"external_ids,omitempty" yaml:"external_ids,omitempty"`
}

// NetworkAddressSetPost used for renaming an Address Set.
//
// swagger:model
//
// API extension: network_address_set.
type NetworkAddressSetPost struct {
	// The new name of the address set
	// Example: "bar"
	Name string `json:"name" yaml:"name"`
}

// NetworkAddressSetPut used for updating an Address Set.
//
// swagger:model
//
// API extension: network_address_set.
type NetworkAddressSetPut struct {
	// Updated addresses
	Addresses []string `json:"addresses" yaml:"addresses"`

	// Updated external_ids
	ExternalIDs map[string]string `json:"external_ids,omitempty" yaml:"external_ids,omitempty"`
}

// NetworkAddressSetsPost used for creating a new Address Set.
//
// swagger:model
//
// API extension: network_address_set.
type NetworkAddressSetsPost struct {
	NetworkAddressSetPut  `yaml:",inline"`
	NetworkAddressSetPost `yaml:",inline"`
}

// Normalise normalises fields in the NetworkAddressSet so that comparisons are consistent.
func (as *NetworkAddressSet) Normalise() {
	as.Name = strings.TrimSpace(as.Name)

	trimmedAddresses := make([]string, 0, len(as.Addresses))
	for _, addr := range as.Addresses {
		trimmedAddresses = append(trimmedAddresses, strings.TrimSpace(addr))
	}
	as.Addresses = trimmedAddresses

	if as.ExternalIDs != nil {
		normalized := make(map[string]string, len(as.ExternalIDs))
		for k, v := range as.ExternalIDs {
			normalized[strings.TrimSpace(k)] = strings.TrimSpace(v)
		}
		as.ExternalIDs = normalized
	}
}

// Writable converts a full NetworkAddressSet struct into a NetworkAddressSetPut struct (filters read-only fields).
func (as *NetworkAddressSet) Writable() NetworkAddressSetPut {
	return NetworkAddressSetPut{
		Addresses:   as.Addresses,
		ExternalIDs: as.ExternalIDs,
	}
}
