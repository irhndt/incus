package address_set

import (
	"fmt"

	"github.com/lxc/incus/v6/internal/server/state"
	"github.com/lxc/incus/v6/shared/logger"
)

// FirewallApplyAddressSetRules applies Address Set rules to the network firewall.
func FirewallApplyAddressSetRules(s *state.State, logger logger.Logger, projectName string, addressSet AddressSetUsage) error {
	// Extract Address Set information.
	setName := addressSet.Name
	addresses := addressSet.Addresses

	// Utilize the nftables driver to create or update nft sets.
	if s.Firewall.String() != "nftables" {
		return fmt.Errorf("Firewall driver nftables not found only supported for now")
	}

	// Apply Address Set to nftables using the updated interface.
	err := s.Firewall.AddressSetToNFTSets(setName, addresses)
	if err != nil {
		return fmt.Errorf("Failed to apply Address Set %q to nftables: %w", setName, err)
	}

	return nil
}

// FirewallClearAddressSetRules removes Address Set rules from the network firewall.
func FirewallClearAddressSetRules(s *state.State, logger logger.Logger, projectName string, setName string) error {
	// Utilize the nftables driver to remove nft sets.
	if s.Firewall.String() != "nftables" {
		return fmt.Errorf("Firewall driver nftables not found only supported for now")
	}

	// Remove Address Set from nftables using AddressSetRemove.
	err := s.Firewall.AddressSetRemove(setName)
	if err != nil {
		return fmt.Errorf("Failed to remove Address Set %q from nftables: %w", setName, err)
	}

	return nil
}