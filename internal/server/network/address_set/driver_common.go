package address_set

import (
	"context"
	"fmt"
	"net"

	incus "github.com/lxc/incus/v6/client"
	"github.com/lxc/incus/v6/internal/server/cluster"
	"github.com/lxc/incus/v6/internal/server/cluster/request"
	"github.com/lxc/incus/v6/internal/server/db"
	"github.com/lxc/incus/v6/internal/server/state"
	localUtil "github.com/lxc/incus/v6/internal/server/util"
	"github.com/lxc/incus/v6/internal/version"
	"github.com/lxc/incus/v6/shared/api"
	"github.com/lxc/incus/v6/shared/logger"
	"github.com/lxc/incus/v6/shared/revert"
)

// common represents a Network Address Set.
type common struct {
	logger      logger.Logger
	state       *state.State
	id          int64
	projectName string
	info        *api.NetworkAddressSet
}

// init initialize internal variables.
func (d *common) init(state *state.State, id int64, projectName string, info *api.NetworkAddressSet) {
	if info == nil {
		d.info = &api.NetworkAddressSet{}
	} else {
		d.info = info
	}

	d.logger = logger.AddContext(logger.Ctx{"project": projectName, "networkAddressSet": d.info.Name})
	d.id = id
	d.projectName = projectName
	d.state = state

	if d.info.Addresses == nil {
		d.info.Addresses = []string{}
	}

	if d.info.ExternalIDs == nil {
		d.info.ExternalIDs = make(map[string]string)
	}
}

// ID returns the Network Address Set ID.
func (d *common) ID() int64 {
	return d.id
}

// Project returns the project name.
func (d *common) Project() string {
	return d.projectName
}

// Info returns copy of internal info for the Network AddressSet.
func (d *common) Info() *api.NetworkAddressSet {
	// Copy internal info to prevent modification externally.
	info := api.NetworkAddressSet{}
	info.Name = d.info.Name
	info.Description = d.info.Description
	info.Addresses = append([]string(nil), d.info.Addresses...)
	info.ExternalIDs = localUtil.CopyConfig(d.info.ExternalIDs)
	info.UsedBy = nil // To indicate its not populated (use Usedby() function to populate).
	info.Project = d.projectName

	return &info
}

// usedBy returns a list of ACLs API endpoints referencing this Address Set.
// If firstOnly is true then search stops at first result.
func (d *common) usedBy(firstOnly bool) ([]string, error) {
	usedBy := []string{}

	// Find all ACLs that reference this address set.
	err := AddressSetUsedBy(d.state, d.projectName, func(aclName string) error {
		uri := fmt.Sprintf("/%s/network-acls/%s", version.APIVersion, aclName)
		if d.projectName != api.ProjectDefaultName {
			uri += fmt.Sprintf("?project=%s", d.projectName)
		}
		usedBy = append(usedBy, uri)

		if firstOnly {
			return db.ErrInstanceListStop
		}

		return nil
	}, d.info.Name)
	if err != nil {
		if err == db.ErrInstanceListStop {
			return usedBy, nil
		}

		return nil, fmt.Errorf("Failed getting address set usage: %w", err)
	}

	return usedBy, nil
}

// UsedBy returns a list of ACL API endpoints referencing this Address Set.
func (d *common) UsedBy() ([]string, error) {
	return d.usedBy(false)
}

// Etag returns the values used for etag generation.
func (d *common) Etag() []any {
	return []any{d.info.Name, d.info.Description, d.info.Addresses, d.info.ExternalIDs}
}

// isUsed returns whether or not the Address Set is used.
func (d *common) isUsed() (bool, error) {
	usedBy, err := d.usedBy(true)
	if err != nil {
		return false, err
	}

	return len(usedBy) > 0, nil
}

// validateName checks name is valid.
func (d *common) validateName(name string) error {
	return ValidName(name)
}

// validateAddresses ensure set is valid.
func (d *common) validateAddresses(addresses []string) error {
	if len(addresses) == 0 {
		// Empty address list is allowed, means no addresses.
		return nil
	}

	var addrType string
	detected := false

	for i, addr := range addresses {
		ip := net.ParseIP(addr)
		if ip != nil {
			// It's an IP address.
			if ip.To4() != nil {
				if !detected {
					addrType = "ipv4"
					detected = true
				} else if addrType != "ipv4" {
					return fmt.Errorf("Mixed address types detected. All addresses must be the same type")
				}
			} else {
				// IPv6 address
				if !detected {
					addrType = "ipv6"
					detected = true
				} else if addrType != "ipv6" {
					return fmt.Errorf("Mixed address types detected. All addresses must be the same type")
				}
			}
		} else {
			// Check MAC
			_, err := net.ParseMAC(addr)
			if err == nil {
				// It's a MAC address
				if !detected {
					addrType = "mac"
					detected = true
				} else if addrType != "mac" {
					return fmt.Errorf("Mixed address types detected. All addresses must be the same type")
				}
			} else {
				return fmt.Errorf("Unsupported address format %q at index %d", addr, i)
			}
		}
	}

	return nil
}

// validateConfig checks the entire config including name and addresses.
func (d *common) validateConfig(config *api.NetworkAddressSetPut) error {
	err := d.validateName(d.info.Name)
	if err != nil {
		return fmt.Errorf("Invalid name: %w", err)
	}

	err = d.validateAddresses(config.Addresses)
	if err != nil {
		return fmt.Errorf("Invalid addresses: %w", err)
	}

	// ExternalIDs can be arbitrary key-value pairs, no validation.
	// Description is free-form text, no validation needed.

	return nil
}

func (d *common) Update(config *api.NetworkAddressSetPut, clientType request.ClientType) error {
	err := d.validateConfig(config)
	if err != nil {
		return err
	}

	revert := revert.New()
	defer revert.Fail()

	if clientType == request.ClientTypeNormal {
		oldConfig := d.info.NetworkAddressSetPut
		err = d.state.DB.Cluster.Transaction(context.TODO(), func(ctx context.Context, tx *db.ClusterTx) error {
			// Update database. Its important this occurs before we attempt to apply to networks
			err := tx.UpdateNetworkAddressSet(ctx, d.projectName, d.info.Name, config)
			return err
		})
		if err != nil {
			return err
		}
		// Apply changes internally and reinitialize.
		d.info.NetworkAddressSetPut = *config
		d.init(d.state, d.id, d.projectName, d.info)
		revert.Add(func() {
			_ = d.state.DB.Cluster.Transaction(context.TODO(), func(ctx context.Context, tx *db.ClusterTx) error {
				return tx.UpdateNetworkAddressSet(ctx, d.projectName, d.info.Name, &oldConfig)
			})

			d.info.NetworkAddressSetPut = oldConfig
			d.init(d.state, d.id, d.projectName, d.info)
		})
	}

	// Get a list of networks that indirectly reference this Address Set via ACLs.
	asNets := map[string]AddressSetUsage{}
	err = AddressSetNetworkUsage(d.state, d.projectName, d.info.Name, d.info.Addresses, asNets)
	if err != nil {
		return fmt.Errorf("Failed getting address set network usage: %w", err)
	}

	// Separate out OVN networks from non-OVN networks for different handling.
	asOVNNets := map[string]AddressSetUsage{}
	for k, v := range asNets {
		if v.Type == "ovn" {
			delete(asNets, k)
			asOVNNets[k] = v
		} else if v.Type != "bridge" {
			return fmt.Errorf("Unsupported network type %q using address set %q", v.Type, d.info.Name)
		}
	}

	// Apply address set changes to non-OVN networks on this member.
	for _, asNet := range asNets {
		err = FirewallApplyAddressSetRules(d.state, d.logger, d.projectName, asNet)
		if err != nil {
			return err
		}
	}

	// If there are affected OVN networks, then apply changes if request type is normal.
	if len(asOVNNets) > 0 && clientType == request.ClientTypeNormal {
		// Check that OVN is available.
		ovnnb, _, err := d.state.OVN()
		if err != nil {
			return err
		}

		// We may need AddressSetNameIDs if we do referencing in OVN. For ACL we had aclNameIDs.
		// For address sets, you might not need name->ID mapping if OVN can handle sets by name directly.
		// If needed, add a similar code block to fetch IDs like done for ACL.

		// Ensure address sets are created or updated in OVN. A function like OVNEnsureAddressSets can be implemented.
		cleanup, err := OVNEnsureAddressSets(d.state, d.logger, ovnnb, d.projectName, asOVNNets, d.info.Name)
		if err != nil {
			return fmt.Errorf("Failed ensuring address set %q is configured in OVN: %w", d.info.Name, err)
		}

		revert.Add(cleanup)

		// If some sets become unused after update, perform cleanup if needed.
		err = OVNAddressSetDeleteIfUnused(d.state, d.logger, ovnnb, d.projectName, d.info.Name)
		if err != nil {
			return fmt.Errorf("Failed removing unused OVN address sets: %w", err)
		}
	}

	// If normal request and asNets is not empty, notify other cluster members.
	if clientType == request.ClientTypeNormal && len(asNets) > 0 {
		notifier, err := cluster.NewNotifier(d.state, d.state.Endpoints.NetworkCert(), d.state.ServerCert(), cluster.NotifyAll)
		if err != nil {
			return err
		}

		err = notifier(func(client incus.InstanceServer) error {
			// Make sure we have a suitable endpoint for updating the address set on other members.
			return client.UseProject(d.projectName).UpdateNetworkAddressSet(d.info.Name, d.info.NetworkAddressSetPut, "")
		})
		if err != nil {
			return err
		}
	}

	revert.Success()
	return nil
}

func (d *common) Rename(newName string) error {
	err := ValidName(newName)
	if err != nil {
		return err
	}

	// Check if name already exists.
	_, err = LoadByName(d.state, d.projectName, newName)
	if err == nil {
		return fmt.Errorf("Address set by that name exists already")
	}

	usedBy, err := d.UsedBy()
	if err != nil {
		return err
	}

	if len(usedBy) > 0 {
		return fmt.Errorf("Cannot rename address set that is in use")
	}

	err = d.state.DB.Cluster.Transaction(context.TODO(), func(ctx context.Context, tx *db.ClusterTx) error {
		return tx.RenameNetworkAddressSet(ctx, d.projectName, d.info.Name, newName)
	})
	if err != nil {
		return err
	}

	d.info.Name = newName

	// TODO: lifecycle.NetworkAddressSetRenamed.Event(d, ...)

	return nil
}

func (d *common) Delete() error {
	usedBy, err := d.UsedBy()
	if err != nil {
		return err
	}

	if len(usedBy) > 0 {
		return fmt.Errorf("Cannot delete address set that is in use")
	}

	return d.state.DB.Cluster.Transaction(context.TODO(), func(ctx context.Context, tx *db.ClusterTx) error {
		return tx.DeleteNetworkAddressSet(ctx, d.projectName, d.info.Name)
	})
}