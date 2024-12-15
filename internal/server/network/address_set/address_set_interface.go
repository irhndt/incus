package address_set

import (
	"github.com/lxc/incus/v6/internal/server/db/cluster/request"
	"github.com/lxc/incus/v6/internal/server/state"
	"github.com/lxc/incus/v6/shared/api"
)

// NetworkAddressSet represents a Network Address Set.
type NetworkAddressSet interface {
	// Initialize.
	init(state *state.State, projectName string, addressSetInfo *api.NetworkAddressSet)

	// Info returns a copy of the address set struct.
	Info() *api.NetworkAddressSet

	// Info
	ID() int64
	Project() string
	Etag() []any
	UsedBy() ([]string, error)

	// Internal validation.
	validateName(name string) error
	validateConfig(config *api.NetworkACLPut) error

	// Modifications.
	Update(config *api.NetworkAddressSetPut, clientType request.ClientType) error
	Rename(newName string) error
	Delete() error
}
