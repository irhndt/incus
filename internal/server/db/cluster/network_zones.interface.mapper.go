//go:build linux && cgo && !agent

package cluster

import (
	"context"
	"github.com/lxc/incus/v6/shared/api"
)

// NetworkZoneGenerated is an interface of generated methods for NetworkZone.
type NetworkZoneGenerated interface {
	// GetNetworkZoneConfig returns all available NetworkZone Configs
	// generator: network_zone GetMany
	GetNetworkZoneConfig(ctx context.Context, db dbtx, networkZoneID int, filters ...ConfigFilter) (map[string]string, error)

	// GetNetworkZones returns the existing Network zones.
	// generator: network_zone GetMany
	GetNetworkZones(ctx context.Context, db dbtx, filters ...ConfigFilter) ([]NetworkZone, error)

	// GetNetworkZone returns the Network zone with the given name in the given project.
	// generator: network_zone GetOne
	GetNetworkZone(ctx context.Context, db dbtx, project string, name string) (*api.NetworkZone, error)

	// CreateNetworkZone creates a new Network zone.
	// generator: network_zone Create
	CreateNetworkZone(ctx context.Context, db dbtx, project string, info *api.NetworkZonesPost) (int64, error)

	// UpdateNetworkZone updates the Network zone with the given ID.
	// generator: network_zone Update
	UpdateNetworkZone(ctx context.Context, db dbtx, id int64, config *api.NetworkZonePut) error

	// DeleteNetworkZone deletes the Network zone.
	// generator: network_zone DeleteOne-by-ID
	DeleteNetworkZone(ctx context.Context, db dbtx, id int64) error

	// Records section

	// GetNetworkZoneRecordConfig returns all available NetworkZone Config
	// generator: network_zone GetMany
	GetNetworkZoneRecordConfig(ctx context.Context, db dbtx, networkZoneRecordID int, filters ...ConfigFilter) (map[string]string, error)

	// GetNetworkZoneRecords returns the network zone record for the given zone and name.
	// generator: network_zone GetMany
	GetNetworkZoneRecords(ctx context.Context, db dbtx, filters ...ConfigFilter) ([]NetworkZoneRecord, error)

	// GetNetworkZoneRecord returns the network zone record for the given zone and name.
	// generator: network_zone GetOne
	GetNetworkZoneRecord(ctx context.Context, db dbtx, zone int64, name string) (*NetworkZoneRecord, error)

	// CreateNetworkZoneRecord creates a new network zone record.
	// generator: network_zone Create
	CreateNetworkZoneRecord(ctx context.Context, db dbtx, zone int64, info api.NetworkZoneRecordsPost) (int64, error)

	// UpdateNetworkZoneRecord updates the network zone record with the given ID.
	// generator: network_zone Update
	UpdateNetworkZoneRecord(ctx context.Context, db dbtx, id int64, config api.NetworkZoneRecordPut) error

	// DeleteNetworkZoneRecord deletes the network zone record.
	// generator: network_zone DeleteOne-by-ID
	DeleteNetworkZoneRecord(ctx context.Context, db dbtx, id int64) error
}
