package address_sets

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/lxc/incus/v6/internal/server/network/ovn"
	"github.com/lxc/incus/v6/internal/server/state"
	"github.com/lxc/incus/v6/shared/logger"
	"github.com/lxc/incus/v6/shared/revert"
)

func OVNEnsureAddressSets(s *state.State, l logger.Logger, client *ovn.NB, projectName string, asNets map[string]AddressSetUsage, addressSetName string) (revert.Hook, error) {
	revert := revert.New()
	defer revert.Fail()

	// Load the address set info.
	addrSet, err := LoadByName(s, projectName, addressSetName)
	if err != nil {
		return nil, fmt.Errorf("Failed loading address set %q: %w", addressSetName, err)
	}

	asInfo := addrSet.Info()

	// Convert addresses to a map to handle IP and MAC separation.
	addressMap := make(map[string]struct{}, len(asInfo.Addresses))
	for _, a := range asInfo.Addresses {
		addressMap[a] = struct{}{}
	}

	// If no addresses, still ensure empty sets exist so ACL match doesn't fail.
	err = client.AddressSetToNFTSets(ovn.OVNAddressSet(asInfo.Name), addressMap)
	if err != nil {
		return nil, fmt.Errorf("Failed applying address set %q to nft sets: %w", asInfo.Name, err)
	}

	ipNets := []net.IPNet{}
	for addr := range addressMap {
		if strings.Count(addr, "/") == 1 {
			_, ipnet, err := net.ParseCIDR(addr)
			if err != nil {
				// Maybe it's a plain IP or MAC, handle gracefully:
				ip := net.ParseIP(addr)
				if ip != nil {
					// Convert to /32 or /128 CIDR
					bits := 32
					if ip.To4() == nil {
						bits = 128
					}
					mask := net.CIDRMask(bits, bits)
					ipnet := net.IPNet{IP: ip, Mask: mask}
					ipNets = append(ipNets, ipnet)
					continue
				}

				// If MAC, skip. address sets can also contain MAC
				// For MAC addresses, handle them differently if needed.
				if _, errMac := net.ParseMAC(addr); errMac == nil {
					// For MAC addresses, we can't just store them as IPNet, skip IPNets for them.
					continue
				}
				return nil, fmt.Errorf("Failed parsing address %q: %w", addr, err)
			}
			ipNets = append(ipNets, *ipnet)
		} else {
			// Might be a single IP address without CIDR.
			ip := net.ParseIP(addr)
			if ip != nil {
				bits := 32
				if ip.To4() == nil {
					bits = 128
				}
				mask := net.CIDRMask(bits, bits)
				ipnet := net.IPNet{IP: ip, Mask: mask}
				ipNets = append(ipNets, ipnet)
			} else {
				// If MAC, skip from IP sets.
				if _, errMac := net.ParseMAC(addr); errMac == nil {
					continue
				}
				return nil, fmt.Errorf("Unsupported address %q in address set", addr)
			}
		}
	}

	// Ensure OVN NB DB address sets exist and are populated.
	err = client.UpdateAddressSetAdd(context.TODO(), ovn.OVNAddressSet(asInfo.Name), ipNets...)
	if err != nil {
		return nil, fmt.Errorf("Failed to ensure address set %q in OVN: %w", asInfo.Name, err)
	}

	revert.Success()
	return nil, nil
}
