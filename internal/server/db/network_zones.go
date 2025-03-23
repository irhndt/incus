//go:build linux && cgo && !agent

package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/lxc/incus/v6/internal/server/db/query"
	"github.com/lxc/incus/v6/internal/version"
	"github.com/lxc/incus/v6/shared/api"
)

// GetNetworkZoneKeys returns a map of key names to keys.
func (c *ClusterTx) GetNetworkZoneKeys(ctx context.Context) (map[string]string, error) {
	q := `SELECT networks_zones.name, networks_zones_config.key, networks_zones_config.value
		FROM networks_zones
		JOIN networks_zones_config ON networks_zones_config.network_zone_id=networks_zones.id
		WHERE networks_zones_config.key LIKE 'peers.%.key'
	`

	secrets := map[string]string{}

	err := query.Scan(ctx, c.tx, q, func(scan func(dest ...any) error) error {
		var name string
		var peer string
		var secret string

		err := scan(&name, &peer, &secret)
		if err != nil {
			return err
		}

		fields := strings.SplitN(peer, ".", 3)
		if len(fields) != 3 {
			// Skip invalid values.
			return nil
		}

		// Format as a valid TSIG secret (encode domain name, key name and make valid FQDN).
		secrets[fmt.Sprintf("%s_%s.", name, fields[1])] = secret

		return nil
	})
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

// networkZoneConfig populates the config map of the Network zone with the given ID.
func networkZoneConfig(ctx context.Context, tx *ClusterTx, id int64, zone *api.NetworkZone) error {
	q := `
		SELECT key, value
		FROM networks_zones_config
		WHERE network_zone_id=?
	`

	zone.Config = make(map[string]string)
	return query.Scan(ctx, tx.Tx(), q, func(scan func(dest ...any) error) error {
		var key, value string

		err := scan(&key, &value)
		if err != nil {
			return err
		}

		_, found := zone.Config[key]
		if found {
			return fmt.Errorf("Duplicate config row found for key %q for network zone ID %d", key, id)
		}

		zone.Config[key] = value

		return nil
	}, id)
}

// networkzoneConfigAdd inserts Network zone config keys.
func networkzoneConfigAdd(tx *sql.Tx, id int64, config map[string]string) error {
	sql := "INSERT INTO networks_zones_config (network_zone_id, key, value) VALUES(?, ?, ?)"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}

	defer func() { _ = stmt.Close() }()

	for k, v := range config {
		if v == "" {
			continue
		}

		_, err = stmt.Exec(id, k, v)
		if err != nil {
			return fmt.Errorf("Failed inserting config: %w", err)
		}
	}

	return nil
}

// networkZoneRecordConfig populates the config map of the network zone record with the given ID.
func networkZoneRecordConfig(ctx context.Context, tx *ClusterTx, id int64, record *api.NetworkZoneRecord) error {
	q := `
		SELECT key, value
		FROM networks_zones_records_config
		WHERE network_zone_record_id=?
	`

	record.Config = make(map[string]string)
	return query.Scan(ctx, tx.Tx(), q, func(scan func(dest ...any) error) error {
		var key, value string

		err := scan(&key, &value)
		if err != nil {
			return err
		}

		_, found := record.Config[key]
		if found {
			return fmt.Errorf("Duplicate config row found for key %q for network zone ID %d", key, id)
		}

		record.Config[key] = value

		return nil
	}, id)
}

// networkzoneConfigAdd inserts Network zone config keys.
func networkZoneRecordConfigAdd(tx *sql.Tx, id int64, config map[string]string) error {
	sql := "INSERT INTO networks_zones_records_config (network_zone_record_id, key, value) VALUES(?, ?, ?)"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}

	defer func() { _ = stmt.Close() }()

	for k, v := range config {
		if v == "" {
			continue
		}

		_, err = stmt.Exec(id, k, v)
		if err != nil {
			return fmt.Errorf("Failed inserting config: %w", err)
		}
	}

	return nil
}

// GetNetworkZoneURIs returns the URIs for the network ACLs with the given project.
func (c *ClusterTx) GetNetworkZoneURIs(ctx context.Context, projectID int, project string) ([]string, error) {
	q := `SELECT networks_zones.name from networks_zones WHERE networks_zones.project_id = ?`

	names, err := query.SelectStrings(ctx, c.tx, q, projectID)
	if err != nil {
		return nil, fmt.Errorf("Unable to get URIs for network zone: %w", err)
	}

	uris := make([]string, len(names))
	for i := range names {
		uris[i] = api.NewURL().Path(version.APIVersion, "network-zones", names[i]).Project(project).String()
	}

	return uris, nil
}
