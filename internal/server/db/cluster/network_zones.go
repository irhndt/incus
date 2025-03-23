//go:build linux && cgo && !agent

package cluster

import (
	"context"
	"database/sql"
	"github.com/lxc/incus/v6/shared/api"
)

// Code generation directives.
//
//generate-database:mapper target networks_zone.mapper.go
//generate-database:mapper reset -i -b "//go:build linux && cgo && !agent"
//

//generate-database:mapper stmt -e network_zone objects
//generate-database:mapper stmt -e network_zone objects-by-Name
//generate-database:mapper stmt -e network_zone objects-by-ID
//generate-database:mapper stmt -e network_zone create struct=NetworkZone table=networks_zones
//generate-database:mapper stmt -e network_zone id
//generate-database:mapper stmt -e network_zone rename
//generate-database:mapper stmt -e network_zone update struct=NetworkZone table=networks_zones
//generate-database:mapper stmt -e network_zone DeleteOne-by-ID table=networks_zones

//generate-database:mapper method -i -e network_zone GetMany references=Config
//generate-database:mapper method -i -e network_zone GetOne struct=NetworkZone
//generate-database:mapper method -i -e network_zone Create references=Config
//generate-database:mapper method -i -e network_zone Update struct=NetworkZone references=Config
//generate-database:mapper method -i -e network_zone DeleteOne-by-ID
//generate-database:mapper method -i -e network_zone Rename

// NetworkZoneRecord
//generate-database:mapper stmt -e network_zone_record objects table=networks_zones_records
//generate-database:mapper stmt -e network_zone_record objects-by-ID table=networks_zones_records
//generate-database:mapper stmt -e network_zone_record objects-by-Zone table=networks_zones_records
//generate-database:mapper stmt -e network_zone_record objects-by-Zone-and-Name table=networks_zones_records
//generate-database:mapper stmt -e network_zone_record create struct=NetworkZoneRecord table=networks_zones_records
//generate-database:mapper stmt -e network_zone_record update struct=NetworkZoneRecord table=networks_zones_records
//generate-database:mapper stmt -e network_zone_record delete-by-ID table=networks_zones_records

//generate-database:mapper method -i -e network_zone_record GetMany references=Config
//generate-database:mapper method -i -e network_zone_record GetOne struct=NetworkZoneRecord
//generate-database:mapper method -i -e network_zone_record Create references=Config
//generate-database:mapper method -i -e network_zone_record Update struct=NetworkZoneRecord references=Config
//generate-database:mapper method -i -e network_zone_record DeleteOne-by-ID

// NetworkZone is a value object holding db-related details about a network zone (DNS).
type NetworkZone struct {
	ID          int    `db:"omit=create,update"`
	ProjectID   int    `db:"omit=create,update"`
	Project     string `db:"primary=yes&join=projects.name"`
	Name        string `db:"primary=yes"`
	Description string `db:"coalesce=''"`
}

// ToAPI converts the DB records to an API record.
func (n *NetworkZone) ToAPI(ctx context.Context, tx *sql.Tx) (*api.NetworkZone, error) {
	// Get the config.
	config, err := GetNetworkZoneConfig(ctx, tx, n.ID)
	if err != nil {
		return nil, err
	}

	// Fill in the struct.
	resp := api.NetworkZone{
		Name: n.Name,
		NetworkZonePut: api.NetworkZonePut{
			Description: n.Description,
			Config:      config,
		},
	}

	return &resp, nil
}

// NetworkZoneFilter specifies potential query parameter fields.
type NetworkZoneFilter struct {
	ID   *int
	Name *string
	Project *string
}

// NetworkZoneRecord is a value object holding db-related details about a DNS record in a network zone.
type NetworkZoneRecord struct {
	ID            int      `db:"omit=create,update"`
	ZoneID        int      `db:"omit=create,update"` // Refers to networks_zones.id
	Zone          string   `db:"primary=yes&join=networks_zones.name"` // Zone name via join
	Name          string   `db:"primary=yes"`                          // DNS record name
	Description   string   `db:"coalesce=''"`
	Entries       []api.NetworkZoneRecordEntry `db:"marshal=json"`
}

func (r *NetworkZoneRecord) ToAPI(ctx context.Context, tx *sql.Tx) (*api.NetworkZoneRecord, error) {
	config, err := GetNetworkZoneRecordConfig(ctx, tx, r.ID)
	if err != nil {
		return nil, err
	}
	record := api.NetworkZoneRecord{
		Name: r.Name,
		NetworkZoneRecordPut: api.NetworkZoneRecordPut{
			Description: 	r.Description,
			Entries: 		r.Entries,
			Config:     	config,
		},
	}
	return &record, nil
}

type NetworkZoneRecordFilter struct {
	ID   	*int
	Name 	*string
	Zone	*string
	ZoneID	*int
}
